package repositories

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	eventomodels "ms-congress/internal/models/evento"
	pedidomodels "ms-congress/internal/models/pedido"
	carrinho "ms-congress/internal/services/carrinho"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PedidoRepository struct{ db *gorm.DB }

var ErrPaymentReferenceConflict = errors.New("referência externa de pagamento já utilizada")
var ErrPaymentStateConflict = errors.New("transição de pagamento inválida")
var ErrPaymentNotFound = errors.New("pagamento não encontrado")
var ErrPaymentDataProtectionNotConfigured = errors.New("chave de proteção dos dados de pagamento não configurada")

func NewPedidoRepository(db *gorm.DB) *PedidoRepository { return &PedidoRepository{db: db} }

func (r *PedidoRepository) BuscarDoUsuario(usuarioID uint, pedidoID uint64) (*pedidomodels.Pedido, error) {
	var pedido pedidomodels.Pedido
	err := r.db.Where("id = ? AND usuario_id = ?", pedidoID, usuarioID).
		Preload("Pagamento").Preload("Itens.EventoIngresso.TipoIngresso").First(&pedido).Error
	return &pedido, err
}

func (r *PedidoRepository) BuscarPorChave(usuarioID uint, key string) (*pedidomodels.Pedido, error) {
	var pedido pedidomodels.Pedido
	err := r.db.Where("usuario_id = ? AND idempotency_key = ?", usuarioID, key).
		Preload("Pagamento").Preload("Itens.EventoIngresso.TipoIngresso").First(&pedido).Error
	return &pedido, err
}

func (r *PedidoRepository) ExpirarSeNecessario(pedidoID uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var pedido pedidomodels.Pedido
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&pedido, pedidoID).Error; err != nil {
			return err
		}
		if pedido.Status != pedidomodels.StatusPedidoPendentePagamento || time.Now().Before(pedido.ExpiraEm) {
			return nil
		}
		if err := tx.Model(&eventomodels.Ingresso{}).Where("pedido_id = ? AND status = ?", pedido.ID, "RESERVADO").
			Updates(map[string]interface{}{"status": "DISPONIVEL", "usuario_id": nil, "pedido_id": nil}).Error; err != nil {
			return err
		}
		if err := tx.Model(&pedidomodels.Pedido{}).Where("id = ?", pedido.ID).Update("status", pedidomodels.StatusPedidoExpirado).Error; err != nil {
			return err
		}
		return tx.Model(&pedidomodels.Pagamento{}).Where("pedido_id = ?", pedido.ID).
			Update("status", pedidomodels.StatusPagamentoRejeitado).Error
	})
}

type PaymentDetails struct {
	CPF         string
	CardLast4   string
	CardBrand   string
	ExpiryMonth int
	ExpiryYear  int
}

func (r *PedidoRepository) RegistrarPagamento(pedidoID uint64, key, externalID string, details PaymentDetails) error {
	encryptedCPF, err := encryptCPF(details.CPF)
	if err != nil {
		return err
	}
	return r.db.Transaction(func(tx *gorm.DB) error {
		var pedido pedidomodels.Pedido
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Preload("Pagamento").First(&pedido, pedidoID).Error; err != nil {
			return err
		}
		if pedido.Pagamento != nil && pedido.Pagamento.ReferenciaExterna != "" &&
			!strings.HasPrefix(pedido.Pagamento.ReferenciaExterna, "pending-") {
			return nil
		}
		if pedido.Status != pedidomodels.StatusPedidoPendentePagamento {
			return ErrPaymentStateConflict
		}
		if pedido.Pagamento == nil {
			return errors.New("pagamento do pedido não encontrado")
		}
		var existing pedidomodels.Pagamento
		if err := tx.Select("id").Where("referencia_externa = ? AND id <> ?", externalID, pedido.Pagamento.ID).
			First(&existing).Error; err == nil {
			return ErrPaymentReferenceConflict
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		err := tx.Model(&pedidomodels.Pagamento{}).Where("pedido_id = ?", pedidoID).
			Updates(map[string]interface{}{
				"idempotency_key": key, "referencia_externa": externalID, "provedor": "configured",
				"cpf_encriptado": encryptedCPF, "cartao_final": details.CardLast4,
				"bandeira_cartao": details.CardBrand, "vencimento_mes": details.ExpiryMonth,
				"vencimento_ano": details.ExpiryYear,
			}).Error
		if err != nil && isPaymentReferenceConflict(err) {
			return ErrPaymentReferenceConflict
		}
		return err
	})
}

func isPaymentReferenceConflict(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505" &&
			(strings.Contains(strings.ToLower(pgErr.ConstraintName), "referencia_externa") ||
				strings.Contains(strings.ToLower(pgErr.ConstraintName), "referência_externa"))
	}
	message := strings.ToLower(err.Error())
	return strings.Contains(message, "referencia_externa") ||
		strings.Contains(message, "referência externa") ||
		strings.Contains(message, "idx_tb_pagamento_referencia_externa")
}

func encryptCPF(cpf string) (string, error) {
	secret := os.Getenv("PAYMENT_DATA_ENCRYPTION_KEY")
	if secret == "" {
		return "", ErrPaymentDataProtectionNotConfigured
	}
	key := sha256.Sum256([]byte(secret))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", fmt.Errorf("cifrar CPF: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("cifrar CPF: %w", err)
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", fmt.Errorf("cifrar CPF: %w", err)
	}
	ciphertext := gcm.Seal(nonce, nonce, []byte(cpf), nil)
	return base64.RawStdEncoding.EncodeToString(ciphertext), nil
}

func (r *PedidoRepository) ConfirmarPagamento(externalID, status string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var payment pedidomodels.Pagamento
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("referencia_externa = ?", externalID).First(&payment).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrPaymentNotFound
			}
			return err
		}
		var pedido pedidomodels.Pedido
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&pedido, payment.PedidoID).Error; err != nil {
			return err
		}
		if err := tx.Preload("Itens").First(&pedido, payment.PedidoID).Error; err != nil {
			return err
		}
		now := time.Now()
		if pedido.Status == pedidomodels.StatusPedidoPendentePagamento && now.After(pedido.ExpiraEm) {
			if err := tx.Model(&eventomodels.Ingresso{}).Where("pedido_id = ? AND status = ?", pedido.ID, "RESERVADO").
				Updates(map[string]interface{}{"status": "DISPONIVEL", "usuario_id": nil, "pedido_id": nil}).Error; err != nil {
				return err
			}
			if err := tx.Model(&pedidomodels.Pedido{}).Where("id = ?", pedido.ID).
				Updates(map[string]interface{}{"status": pedidomodels.StatusPedidoExpirado, "cancelado_em": now}).Error; err != nil {
				return err
			}
			if err := tx.Model(&pedidomodels.Pagamento{}).Where("id = ?", payment.ID).
				Update("status", pedidomodels.StatusPagamentoRejeitado).Error; err != nil {
				return err
			}
			return ErrPaymentStateConflict
		}
		if status == pedidomodels.StatusPagamentoAprovado &&
			(payment.Status == pedidomodels.StatusPagamentoRejeitado ||
				pedido.Status == pedidomodels.StatusPedidoCancelado ||
				pedido.Status == pedidomodels.StatusPedidoExpirado) {
			return ErrPaymentStateConflict
		}
		if status == pedidomodels.StatusPagamentoRejeitado &&
			(payment.Status == pedidomodels.StatusPagamentoAprovado ||
				pedido.Status == pedidomodels.StatusPedidoPago) {
			return ErrPaymentStateConflict
		}
		if status == pedidomodels.StatusPagamentoRejeitado &&
			payment.Status == pedidomodels.StatusPagamentoRejeitado &&
			(pedido.Status == pedidomodels.StatusPedidoCancelado ||
				pedido.Status == pedidomodels.StatusPedidoExpirado) {
			return nil
		}
		if status == pedidomodels.StatusPagamentoAprovado {
			var reserved []eventomodels.Ingresso
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("pedido_id = ? AND status = ?", pedido.ID, "RESERVADO").Find(&reserved).Error; err != nil {
				return err
			}
			// A callback aprovado pode ser repetido depois de uma emissão
			// parcial ou de uma versão anterior da aplicação. Só encerramos
			// quando não há reservas pendentes para reconciliar.
			if len(reserved) == 0 && payment.Status == pedidomodels.StatusPagamentoAprovado && pedido.Status == pedidomodels.StatusPedidoPago {
				var issued int64
				if err := tx.Model(&eventomodels.Ingresso{}).
					Where("pedido_id = ? AND status = ?", pedido.ID, "VENDIDO").
					Count(&issued).Error; err != nil {
					return err
				}
				expected := 0
				for _, item := range pedido.Itens {
					expected += item.Quantidade
				}
				if int(issued) == expected && expected > 0 {
					return nil
				}
				return emitMissingPaidOrderTickets(tx, pedido, expected-int(issued), now)
			}
			if err := ensureOrderTicketCount(tx, pedido, len(reserved)); err != nil {
				return err
			}
			for _, ingresso := range reserved {
				credentialCode, err := newCredentialCode()
				if err != nil {
					return err
				}

				if err := tx.Model(&eventomodels.Ingresso{}).Where("id = ? AND status = ?", ingresso.ID, "RESERVADO").
					Updates(map[string]interface{}{"status": "VENDIDO", "usuario_id": pedido.UsuarioID, "vendido_em": now, "codigo": credentialCode}).Error; err != nil {
					return err
				}
			}

			if err := tx.Model(&pedidomodels.Pedido{}).Where("id = ?", pedido.ID).
				Updates(map[string]interface{}{"status": pedidomodels.StatusPedidoPago, "pago_em": now}).Error; err != nil {
				return err
			}
			return tx.Model(&pedidomodels.Pagamento{}).Where("id = ?", payment.ID).
				Updates(map[string]interface{}{"status": status, "aprovado_em": now}).Error
		}
		if err := tx.Model(&eventomodels.Ingresso{}).Where("pedido_id = ? AND status = ?", pedido.ID, "RESERVADO").
			Updates(map[string]interface{}{"status": "DISPONIVEL", "usuario_id": nil, "pedido_id": nil}).Error; err != nil {
			return err
		}
		if err := tx.Model(&pedidomodels.Pedido{}).Where("id = ?", pedido.ID).
			Updates(map[string]interface{}{"status": pedidomodels.StatusPedidoCancelado, "cancelado_em": now}).Error; err != nil {
			return err
		}
		return tx.Model(&pedidomodels.Pagamento{}).Where("id = ?", payment.ID).Update("status", status).Error
	})
}

func newCredentialCode() (string, error) {
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", fmt.Errorf("gerar credencial digital: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(randomBytes), nil
}

func emitMissingPaidOrderTickets(tx *gorm.DB, pedido pedidomodels.Pedido, missing int, now time.Time) error {
	if missing <= 0 {
		return errors.New("pedido pago sem ingressos emitidos")
	}
	for _, item := range pedido.Itens {
		var issued int64
		if err := tx.Model(&eventomodels.Ingresso{}).
			Where("pedido_id = ? AND evento_ingresso_id = ? AND status IN ?", pedido.ID, item.EventoIngressoID, []string{"VENDIDO", "RESERVADO"}).
			Count(&issued).Error; err != nil {
			return err
		}
		quantity := item.Quantidade - int(issued)
		if quantity <= 0 {
			continue
		}
		var available []eventomodels.Ingresso
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("evento_ingresso_id = ? AND status = ?", item.EventoIngressoID, "DISPONIVEL").
			Limit(quantity).Find(&available).Error; err != nil {
			return err
		}
		if len(available) != quantity {
			return fmt.Errorf("estoque insuficiente para reconciliar o pedido %d", pedido.ID)
		}
		for _, ingresso := range available {
			code, err := newCredentialCode()
			if err != nil {
				return err
			}
			if err := tx.Model(&eventomodels.Ingresso{}).Where("id = ? AND status = ?", ingresso.ID, "DISPONIVEL").
				Updates(map[string]interface{}{"status": "VENDIDO", "usuario_id": pedido.UsuarioID, "pedido_id": pedido.ID, "vendido_em": now, "codigo": code}).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func ensureOrderTicketCount(tx *gorm.DB, pedido pedidomodels.Pedido, reservedCount int) error {
	expected := 0
	for _, item := range pedido.Itens {
		expected += item.Quantidade
	}
	if expected == 0 || reservedCount != expected {
		return fmt.Errorf("quantidade de ingressos reservados inconsistente: esperado %d, encontrado %d", expected, reservedCount)
	}
	return nil
}

func (r *PedidoRepository) Criar(usuarioID uint, key string, entries []carrinho.CartEntry) (*pedidomodels.Pedido, error) {
	var pedido pedidomodels.Pedido
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("idempotency_key = ?", key).
			Preload("Itens.EventoIngresso.TipoIngresso").Preload("Pagamento").
			First(&pedido).Error; err == nil {
			if pedido.UsuarioID != usuarioID {
				return errors.New("chave de idempotência já utilizada")
			}
			return nil
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		pedido = pedidomodels.Pedido{
			UsuarioID: usuarioID, Status: pedidomodels.StatusPedidoPendentePagamento,
			IdempotencyKey: key, ExpiraEm: time.Now().Add(30 * time.Minute),
		}
		reservas := make([]eventomodels.Ingresso, 0)
		for _, entry := range entries {
			var lote eventomodels.EventoIngresso
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				Preload("Evento").Preload("TipoIngresso").
				First(&lote, entry.EventoIngressoID).Error; err != nil {
				return err
			}
			now := time.Now()
			if !lote.Evento.Publicado || lote.Evento.DataFim.Before(now) ||
				lote.DataInicioVenda.After(now) || lote.DataFimVenda.Before(now) {
				return fmt.Errorf("ingresso fora do período de venda")
			}
			var ingressos []eventomodels.Ingresso
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("evento_ingresso_id = ? AND status = ?", lote.ID, "DISPONIVEL").
				Limit(entry.Quantidade).Find(&ingressos).Error; err != nil {
				return err
			}
			if len(ingressos) < entry.Quantidade {
				return fmt.Errorf("quantidade solicitada indisponível")
			}
			subtotal := lote.Preco * float64(entry.Quantidade)
			pedido.Itens = append(pedido.Itens, pedidomodels.PedidoItem{
				EventoIngressoID: lote.ID, Quantidade: entry.Quantidade,
				PrecoUnitario: lote.Preco, Subtotal: subtotal,
			})
			pedido.Subtotal += subtotal
			reservas = append(reservas, ingressos...)
		}
		pedido.Total = pedido.Subtotal
		// Os itens são persistidos explicitamente logo abaixo, depois que o
		// pedido recebe seu ID. Evita que o GORM salve a associação
		// automaticamente e depois tente inserir os mesmos itens novamente.
		if err := tx.Omit("Itens").Create(&pedido).Error; err != nil {
			return err
		}
		for index := range pedido.Itens {
			pedido.Itens[index].PedidoID = pedido.ID
		}
		if err := tx.Create(&pedido.Itens).Error; err != nil {
			return err
		}
		for index := range reservas {
			reservas[index].Status = "RESERVADO"
			reservas[index].PedidoID = &pedido.ID
			if err := tx.Model(&eventomodels.Ingresso{}).Where("id = ?", reservas[index].ID).
				Updates(map[string]interface{}{"status": reservas[index].Status, "pedido_id": pedido.ID}).Error; err != nil {
				return err
			}
		}
		pedido.Pagamento = &pedidomodels.Pagamento{
			PedidoID: pedido.ID, Status: pedidomodels.StatusPagamentoPendente,
			Valor: pedido.Total, IdempotencyKey: key,
		}
		pendingReference, err := newPaymentReference("pending")
		if err != nil {
			return err
		}
		pedido.Pagamento.ReferenciaExterna = pendingReference
		if err := tx.Create(pedido.Pagamento).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &pedido, nil
}

func newPaymentReference(prefix string) (string, error) {
	randomBytes := make([]byte, 24)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", fmt.Errorf("gerar referência de pagamento: %w", err)
	}
	return fmt.Sprintf("%s-%s", prefix, base64.RawURLEncoding.EncodeToString(randomBytes)), nil
}
