package models

import (
	"time"

	eventomodels "ms-congress/internal/models/evento"
)

const (
	StatusPedidoPendentePagamento = "PENDING_PAYMENT"
	StatusPedidoPago              = "PAID"
	StatusPedidoCancelado         = "CANCELED"
	StatusPedidoExpirado          = "EXPIRED"
)

const (
	StatusPagamentoPendente    = "PENDING"
	StatusPagamentoAprovado    = "APPROVED"
	StatusPagamentoRejeitado   = "REJECTED"
	StatusPagamentoReembolsado = "REFUNDED"
)

type Pedido struct {
	ID             uint64    `gorm:"primaryKey"`
	UsuarioID      uint      `gorm:"not null;index"`
	Status         string    `gorm:"not null;index"`
	Subtotal       float64   `gorm:"not null"`
	Total          float64   `gorm:"not null"`
	IdempotencyKey string    `gorm:"not null;uniqueIndex"`
	ExpiraEm       time.Time `gorm:"not null;index"`
	PagoEm         *time.Time
	CanceladoEm    *time.Time
	CriadoEm       time.Time `gorm:"not null;autoCreateTime"`
	AtualizadoEm   time.Time `gorm:"not null;autoUpdateTime"`

	Itens     []PedidoItem `gorm:"foreignKey:PedidoID"`
	Pagamento *Pagamento   `gorm:"foreignKey:PedidoID"`
}

func (*Pedido) TableName() string { return "tb_pedido" }

type PedidoItem struct {
	ID               uint64                      `gorm:"primaryKey"`
	PedidoID         uint64                      `gorm:"not null;index;uniqueIndex:uk_pedido_lote"`
	EventoIngressoID uint64                      `gorm:"not null;index;uniqueIndex:uk_pedido_lote"`
	Quantidade       int                         `gorm:"not null"`
	PrecoUnitario    float64                     `gorm:"not null"`
	Subtotal         float64                     `gorm:"not null"`
	CriadoEm         time.Time                   `gorm:"not null;autoCreateTime"`
	EventoIngresso   eventomodels.EventoIngresso `gorm:"foreignKey:EventoIngressoID"`
}

func (*PedidoItem) TableName() string { return "tb_pedido_item" }

type Pagamento struct {
	ID                uint64  `gorm:"primaryKey"`
	PedidoID          uint64  `gorm:"not null;uniqueIndex"`
	Status            string  `gorm:"not null;index"`
	Valor             float64 `gorm:"not null"`
	Provedor          string
	ReferenciaExterna string `gorm:"uniqueIndex"`
	IdempotencyKey    string `gorm:"uniqueIndex"`
	AprovadoEm        *time.Time
	ReembolsadoEm     *time.Time
	CriadoEm          time.Time `gorm:"not null;autoCreateTime"`
	AtualizadoEm      time.Time `gorm:"not null;autoUpdateTime"`
}

func (*Pagamento) TableName() string { return "tb_pagamento" }
