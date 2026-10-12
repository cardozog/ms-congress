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
	err := r.db.Joins("JOIN tb_pedido ON tb_pedido.id = ingressos.pedido_id").
		Joins("JOIN tb_pagamento ON tb_pagamento.pedido_id = tb_pedido.id").
		Where("ingressos.usuario_id = ? AND ingressos.status = ? AND tb_pedido.status = ? AND tb_pagamento.status = ?", usuarioID, "VENDIDO", pedidomodels.StatusPedidoPago, pedidomodels.StatusPagamentoAprovado).
		Preload("EventoIngresso.Evento.Organizador").Preload("EventoIngresso.TipoIngresso").
		Order("ingressos.id ASC").Find(&ingressos).Error
	return ingressos, err
}

func (r *MeusEventosRepository) ListarPorEvento(usuarioID uint, eventoID uint64) ([]eventomodels.Ingresso, error) {
	var ingressos []eventomodels.Ingresso
	err := r.db.Joins("JOIN tb_pedido ON tb_pedido.id = ingressos.pedido_id").
		Joins("JOIN evento_ingressos ON evento_ingressos.id = ingressos.evento_ingresso_id").
		Joins("JOIN tb_pagamento ON tb_pagamento.pedido_id = tb_pedido.id").
		Where("ingressos.usuario_id = ? AND ingressos.status = ? AND tb_pedido.status = ? AND tb_pagamento.status = ? AND evento_ingressos.evento_id = ?", usuarioID, "VENDIDO", pedidomodels.StatusPedidoPago, pedidomodels.StatusPagamentoAprovado, eventoID).
		Preload("EventoIngresso.Evento.Organizador").Preload("EventoIngresso.TipoIngresso").
		Order("ingressos.id ASC").Find(&ingressos).Error
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
