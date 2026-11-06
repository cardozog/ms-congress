package repositories

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	eventomodels "ms-congress/internal/models/evento"
	pedidomodels "ms-congress/internal/models/pedido"
	"time"
)

var (
	ErrIngressoNaoEncontrado = errors.New("ingresso não encontrado")
	ErrIngressoJaUtilizado   = errors.New("ingresso já utilizado")
	ErrCheckinForaPeriodo    = errors.New("check-in fora do período")
)

type MeusEventosRepository struct{ db *gorm.DB }

func NewMeusEventosRepository(db *gorm.DB) *MeusEventosRepository {
	return &MeusEventosRepository{db: db}
}

func (r *MeusEventosRepository) Listar(usuarioID uint) ([]eventomodels.Ingresso, error) {
	var ingressos []eventomodels.Ingresso
	err := r.db.Table("ingressos AS i").
		Joins("JOIN tb_pedido AS p ON p.id = i.pedido_id").
		Joins("JOIN tb_pagamento AS pg ON pg.pedido_id = p.id").
		Where("p.usuario_id = ? AND i.status IN ? AND p.status = ? AND pg.status = ?", usuarioID, []string{eventomodels.StatusIngressoVendido, eventomodels.StatusIngressoCheckin}, pedidomodels.StatusPedidoPago, pedidomodels.StatusPagamentoAprovado).
		Preload("EventoIngresso.Evento.Organizador").Preload("EventoIngresso.TipoIngresso").
		Order("i.id ASC").Find(&ingressos).Error
	return ingressos, err
}

func (r *MeusEventosRepository) ListarPorEvento(usuarioID uint, eventoID uint64) ([]eventomodels.Ingresso, error) {
	var ingressos []eventomodels.Ingresso
	err := r.db.Table("ingressos AS i").
		Joins("JOIN tb_pedido AS p ON p.id = i.pedido_id").
		Joins("JOIN evento_ingressos AS ei ON ei.id = i.evento_ingresso_id").
		Joins("JOIN tb_pagamento AS pg ON pg.pedido_id = p.id").
		Where("p.usuario_id = ? AND i.status IN ? AND p.status = ? AND pg.status = ? AND ei.evento_id = ?", usuarioID, []string{eventomodels.StatusIngressoVendido, eventomodels.StatusIngressoCheckin}, pedidomodels.StatusPedidoPago, pedidomodels.StatusPagamentoAprovado, eventoID).
		Preload("EventoIngresso.Evento.Organizador").Preload("EventoIngresso.TipoIngresso").
		Order("i.id ASC").Find(&ingressos).Error
	return ingressos, err
}

func (r *MeusEventosRepository) CheckIn(usuarioID uint, eventoID, ingressoID uint64, now time.Time) (*eventomodels.Ingresso, error) {
	var ingresso eventomodels.Ingresso
	err := r.db.Transaction(func(tx *gorm.DB) error {
		query := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Joins("JOIN evento_ingressos AS ei ON ei.id = ingressos.evento_ingresso_id").
			Joins("JOIN eventos AS e ON e.id = ei.evento_id").
			Joins("JOIN tb_pedido AS p ON p.id = ingressos.pedido_id").
			Joins("JOIN tb_pagamento AS pg ON pg.pedido_id = p.id").
			Where("ingressos.id = ? AND ei.evento_id = ? AND p.usuario_id = ? AND p.status = ? AND pg.status = ?", ingressoID, eventoID, usuarioID, pedidomodels.StatusPedidoPago, pedidomodels.StatusPagamentoAprovado).
			Preload("EventoIngresso.Evento").Preload("EventoIngresso.TipoIngresso").
			First(&ingresso)
		if errors.Is(query.Error, gorm.ErrRecordNotFound) {
			return ErrIngressoNaoEncontrado
		}
		if query.Error != nil {
			return query.Error
		}
		if ingresso.Status == eventomodels.StatusIngressoCheckin {
			return ErrIngressoJaUtilizado
		}
		if ingresso.Status != eventomodels.StatusIngressoVendido {
			return errors.New("ingresso não disponível para check-in")
		}

		var evento eventomodels.Evento
		if err := tx.Model(&eventomodels.Evento{}).Where("id = ?", eventoID).First(&evento).Error; err != nil {
			return err
		}
		if now.Before(evento.DataInicio) || now.After(evento.DataFim) {
			return ErrCheckinForaPeriodo
		}
		if err := tx.Model(&eventomodels.Ingresso{}).
			Where("id = ? AND status = ?", ingressoID, eventomodels.StatusIngressoVendido).
			Updates(map[string]interface{}{"status": eventomodels.StatusIngressoCheckin, "checkin_em": now}).Error; err != nil {
			return err
		}
		ingresso.Status = eventomodels.StatusIngressoCheckin
		ingresso.CheckinEm = &now
		return nil
	})
	return &ingresso, err
}

func ClassificarEvento(evento eventomodels.Evento, now time.Time) string {
	if evento.DataInicio.After(now) {
		return "UPCOMING"
	}
	if evento.DataFim.Before(now) {
		return "ENDED"
	}
	return "ONGOING"
}
