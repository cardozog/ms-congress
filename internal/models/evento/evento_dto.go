package models

import "time"

type EventoPorOrganizadorDTO struct {
	ID            uint64    `json:"id"`
	Nome          string    `json:"nome"`
	Descricao     string    `json:"descricao"`
	DataInicio    time.Time `json:"dataInicio"`
	DataFim       time.Time `json:"dataFim"`
	OrganizadorID uint64    `json:"organizadorID"`
	EventoLogo    string    `json:"eventoLogo"`
}

func NovoEventoPorOrganizadorDTO(evento Evento) EventoPorOrganizadorDTO {
	return EventoPorOrganizadorDTO{
		ID:            evento.ID,
		Nome:          evento.Nome,
		Descricao:     evento.Descricao,
		DataInicio:    evento.DataInicio,
		DataFim:       evento.DataFim,
		OrganizadorID: evento.OrganizadorID,
		EventoLogo:    evento.EventoLogo,
	}
}
