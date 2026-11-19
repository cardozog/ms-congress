package repositories

import (
	eventomodels "ms-congress/internal/models/evento"
	pedidomodels "ms-congress/internal/models/pedido"
	"time"

	"gorm.io/gorm"
)

type EventoRelatorio struct {
	ID                   uint64    `json:"id"`
	Nome                 string    `json:"nome"`
	DataInicio           time.Time `json:"dataInicio"`
	DataFim              time.Time `json:"dataFim"`
	IngressosVendidos    int64     `json:"ingressosVendidos"`
	Presencas            int64     `json:"presencas"`
	ValorArrecadado      float64   `json:"valorArrecadado"`
}

type RelatorioRepository struct{ db *gorm.DB }

func NewRelatorioRepository(db *gorm.DB) *RelatorioRepository {
	return &RelatorioRepository{db: db}
}

func (r *RelatorioRepository) ListarEventos(organizadorUsuarioID uint, organizadorID uint64, dataDe, dataAte *time.Time) ([]EventoRelatorio, error) {
	ticketTotals := r.db.Table("ingressos AS i").
		Select("ei.evento_id, COUNT(i.id) AS ingressos_vendidos, SUM(CASE WHEN i.status = ? THEN 1 ELSE 0 END) AS presencas", eventomodels.StatusIngressoCheckin).
		Joins("JOIN evento_ingressos AS ei ON ei.id = i.evento_ingresso_id").
		Joins("JOIN tb_pedido AS p ON p.id = i.pedido_id").
		Joins("JOIN tb_pagamento AS pg ON pg.pedido_id = p.id").
		Where("i.status IN ? AND p.status = ? AND pg.status = ?",
			[]string{eventomodels.StatusIngressoVendido, eventomodels.StatusIngressoCheckin},
			pedidomodels.StatusPedidoPago, pedidomodels.StatusPagamentoAprovado).
		Group("ei.evento_id")

	revenueTotals := r.db.Table("tb_pedido_item AS pi").
		Select("ei.evento_id, SUM(pi.subtotal) AS valor_arrecadado").
		Joins("JOIN evento_ingressos AS ei ON ei.id = pi.evento_ingresso_id").
		Joins("JOIN tb_pedido AS p ON p.id = pi.pedido_id").
		Joins("JOIN tb_pagamento AS pg ON pg.pedido_id = p.id").
		Where("p.status = ? AND pg.status = ?",
			pedidomodels.StatusPedidoPago, pedidomodels.StatusPagamentoAprovado).
		Group("ei.evento_id")

	relatorios := make([]EventoRelatorio, 0)
	query := r.db.Table("eventos AS e").
		Select("e.id, e.nome, e.data_inicio, e.data_fim, "+
			"COALESCE(ticket_totals.ingressos_vendidos, 0) AS ingressos_vendidos, "+
			"COALESCE(ticket_totals.presencas, 0) AS presencas, "+
			"COALESCE(revenue_totals.valor_arrecadado, 0) AS valor_arrecadado").
		Joins("JOIN tb_pessoa_juridica AS pj ON pj.id = e.organizador_id").
		Joins("LEFT JOIN (?) AS ticket_totals ON ticket_totals.evento_id = e.id", ticketTotals).
		Joins("LEFT JOIN (?) AS revenue_totals ON revenue_totals.evento_id = e.id", revenueTotals).
		Where("pj.usuario_id = ? AND pj.id = ?", organizadorUsuarioID, organizadorID)
	if dataDe != nil {
		query = query.Where("e.data_fim >= ?", *dataDe)
	}
	if dataAte != nil {
		query = query.Where("e.data_inicio < ?", dataAte.AddDate(0, 0, 1))
	}
	err := query.Order("e.data_inicio DESC, e.id DESC").Scan(&relatorios).Error
	return relatorios, err
}
