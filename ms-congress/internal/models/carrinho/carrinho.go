package models

import (
	eventomodels "ms-congress/internal/models/evento"
	"time"
)

type Carrinho struct {
	ID           uint64         `gorm:"primaryKey"`
	UsuarioID    uint           `gorm:"not null;uniqueIndex"`
	CriadoEm     time.Time      `gorm:"not null;autoCreateTime"`
	AtualizadoEm time.Time      `gorm:"not null;autoUpdateTime"`
	Itens        []CarrinhoItem `gorm:"foreignKey:CarrinhoID"`
}

func (*Carrinho) TableName() string { return "tb_carrinho" }

type CarrinhoItem struct {
	ID               uint64                      `gorm:"primaryKey"`
	CarrinhoID       uint64                      `gorm:"not null;index;uniqueIndex:uk_carrinho_lote"`
	EventoIngressoID uint64                      `gorm:"not null;index;uniqueIndex:uk_carrinho_lote"`
	Quantidade       int                         `gorm:"not null"`
	PrecoUnitario    float64                     `gorm:"not null"`
	CriadoEm         time.Time                   `gorm:"not null;autoCreateTime"`
	AtualizadoEm     time.Time                   `gorm:"not null;autoUpdateTime"`
	EventoIngresso   eventomodels.EventoIngresso `gorm:"foreignKey:EventoIngressoID"`
}

func (*CarrinhoItem) TableName() string { return "tb_carrinho_item" }
