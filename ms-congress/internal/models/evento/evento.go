package models

import (
	"time"

	usuariomodels "ms-congress/internal/models/usuario"
)

type Evento struct {
	ID uint64 `gorm:"primaryKey"`

	Nome       string    `gorm:"not null"`
	Descricao  string    `gorm:"not null"`
	DataInicio time.Time `gorm:"not null"`
	DataFim    time.Time `gorm:"not null"`
	Publicado  bool      `gorm:"not null;default:false"`

	OrganizadorID uint64                       `gorm:"not null;index"`
	Organizador   usuariomodels.PessoaJuridica `gorm:"foreignKey:OrganizadorID"`

	Logradouro string `gorm:"not null"`
	Cidade     string `gorm:"not null"`
	Estado     string `gorm:"not null"`
	CEP        string `gorm:"not null"`
	EventoLogo string `gorm:"null"`

	Ingressos []EventoIngresso `gorm:"foreignKey:EventoID"`
}

type TipoIngresso struct {
	ID   uint64 `gorm:"primaryKey"`
	Nome string `gorm:"not null;unique"`
}

type EventoIngresso struct {
	ID uint64 `gorm:"primaryKey"`

	EventoID       uint64 `gorm:"not null;index;uniqueIndex:uk_evento_tipo"`
	TipoIngressoID uint64 `gorm:"not null;index;uniqueIndex:uk_evento_tipo"`

	Preco      float64 `gorm:"not null"`
	Quantidade int     `gorm:"not null"`

	DataInicioVenda time.Time `gorm:"not null"`
	DataFimVenda    time.Time `gorm:"not null"`

	Evento       Evento       `gorm:"foreignKey:EventoID"`
	TipoIngresso TipoIngresso `gorm:"foreignKey:TipoIngressoID"`

	Ingressos []Ingresso `gorm:"foreignKey:EventoIngressoID"`
}

type Ingresso struct {
	ID uint64 `gorm:"primaryKey"`

	EventoIngressoID uint64 `gorm:"not null;index"`

	Codigo string `gorm:"not null;uniqueIndex"`
	Status string `gorm:"not null"`

	EventoIngresso EventoIngresso `gorm:"foreignKey:EventoIngressoID"`
}
