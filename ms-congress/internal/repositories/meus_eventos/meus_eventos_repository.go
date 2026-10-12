package repositories

import (
	"gorm.io/gorm"
	eventomodels "ms-congress/internal/models/evento"
	pedidomodels "ms-congress/internal/models/pedido"
	"time"
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
		Where("p.usuario_id = ? AND i.status = ? AND p.status = ? AND pg.status = ?", usuarioID, "VENDIDO", pedidomodels.StatusPedidoPago, pedidomodels.StatusPagamentoAprovado).
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
		Where("p.usuario_id = ? AND i.status = ? AND p.status = ? AND pg.status = ? AND ei.evento_id = ?", usuarioID, "VENDIDO", pedidomodels.StatusPedidoPago, pedidomodels.StatusPagamentoAprovado, eventoID).
		Preload("EventoIngresso.Evento.Organizador").Preload("EventoIngresso.TipoIngresso").
		Order("i.id ASC").Find(&ingressos).Error
	return ingressos, err
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
