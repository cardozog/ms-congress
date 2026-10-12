package repositories

import (
	"errors"
	"fmt"
	"time"

	eventomodels "ms-congress/internal/models/evento"
	pedidomodels "ms-congress/internal/models/pedido"
	carrinho "ms-congress/internal/services/carrinho"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PedidoRepository struct{ db *gorm.DB }

func NewPedidoRepository(db *gorm.DB) *PedidoRepository { return &PedidoRepository{db: db} }

func (r *PedidoRepository) BuscarDoUsuario(usuarioID uint, pedidoID uint64) (*pedidomodels.Pedido, error) {
	var pedido pedidomodels.Pedido
	err := r.db.Where("id = ? AND usuario_id = ?", pedidoID, usuarioID).
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

func (r *PedidoRepository) RegistrarPagamento(pedidoID uint64, key, externalID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var pedido pedidomodels.Pedido
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Preload("Pagamento").First(&pedido, pedidoID).Error; err != nil {
			return err
		}
		if pedido.Pagamento != nil && pedido.Pagamento.ReferenciaExterna != "" {
			return nil
		}
		return tx.Model(&pedidomodels.Pagamento{}).Where("pedido_id = ?", pedidoID).
			Updates(map[string]interface{}{"idempotency_key": key, "referencia_externa": externalID, "provedor": "configured"}).Error
	})
}

func (r *PedidoRepository) ConfirmarPagamento(externalID, status string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var payment pedidomodels.Pagamento
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("referencia_externa = ?", externalID).First(&payment).Error; err != nil {
			return err
		}
		var pedido pedidomodels.Pedido
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&pedido, payment.PedidoID).Error; err != nil {
			return err
		}
		if payment.Status == pedidomodels.StatusPagamentoAprovado || pedido.Status == pedidomodels.StatusPedidoPago {
			return nil
		}
		now := time.Now()
		if status == pedidomodels.StatusPagamentoAprovado {
			var reserved []eventomodels.Ingresso
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("pedido_id = ? AND status = ?", pedido.ID, "RESERVADO").Find(&reserved).Error; err != nil {
				return err
			}
			if len(reserved) == 0 {
				return errors.New("nenhum ingresso reservado para o pedido")
			}
			for _, ingresso := range reserved {
				if err := tx.Model(&eventomodels.Ingresso{}).Where("id = ? AND status = ?", ingresso.ID, "RESERVADO").
					Updates(map[string]interface{}{"status": "VENDIDO", "usuario_id": pedido.UsuarioID, "vendido_em": now}).Error; err != nil {
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
			Update("status", pedidomodels.StatusPedidoCancelado).Error; err != nil {
			return err
		}
		return tx.Model(&pedidomodels.Pagamento{}).Where("id = ?", payment.ID).Update("status", status).Error
	})
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
