package models

import "time"

type EventoOrganizadorDTO struct {
	ID           uint   `json:"id"`
	UsuarioID    uint   `json:"usuarioID"`
	RazaoSocial  string `json:"razaoSocial"`
	NomeFantasia string `json:"nomeFantasia"`
	CNPJ         string `json:"cnpj"`
}

type TipoIngressoDTO struct {
	ID   uint64 `json:"id"`
	Nome string `json:"nome"`
}

type EventoIngressoDTO struct {
	ID              uint64          `json:"id"`
	EventoID        uint64          `json:"eventoID"`
	TipoIngressoID  uint64          `json:"tipoIngressoID"`
	Preco           float64         `json:"preco"`
	Quantidade      int             `json:"quantidade"`
	Disponivel      int             `json:"disponivel"`
	DataInicioVenda time.Time       `json:"dataInicioVenda"`
	DataFimVenda    time.Time       `json:"dataFimVenda"`
	TipoIngresso    TipoIngressoDTO `json:"tipoIngresso"`
}

type EventoDetalhadoDTO struct {
	ID            uint64               `json:"id"`
	Nome          string               `json:"nome"`
	Descricao     string               `json:"descricao"`
	DataInicio    time.Time            `json:"dataInicio"`
	DataFim       time.Time            `json:"dataFim"`
	OrganizadorID uint64               `json:"organizadorID"`
	Organizador   EventoOrganizadorDTO `json:"organizador"`
	Logradouro    string               `json:"logradouro"`
	Cidade        string               `json:"cidade"`
	Estado        string               `json:"estado"`
	CEP           string               `json:"cep"`
	EventoLogo    string               `json:"eventoLogo"`
	Ingressos     []EventoIngressoDTO  `json:"ingressos"`
}

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

func NovoEventoDetalhadoDTO(evento Evento) EventoDetalhadoDTO {
	ingressos := make([]EventoIngressoDTO, 0, len(evento.Ingressos))
	for _, ingresso := range evento.Ingressos {
		ingressos = append(ingressos, EventoIngressoDTO{
			ID: ingresso.ID, EventoID: ingresso.EventoID, TipoIngressoID: ingresso.TipoIngressoID,
			Preco: ingresso.Preco, Quantidade: ingresso.Quantidade, Disponivel: ingresso.Disponivel,
			DataInicioVenda: ingresso.DataInicioVenda, DataFimVenda: ingresso.DataFimVenda,
			TipoIngresso: TipoIngressoDTO{ID: ingresso.TipoIngresso.ID, Nome: ingresso.TipoIngresso.Nome},
		})
	}

	return EventoDetalhadoDTO{
		ID: evento.ID, Nome: evento.Nome, Descricao: evento.Descricao,
		DataInicio: evento.DataInicio, DataFim: evento.DataFim, OrganizadorID: evento.OrganizadorID,
		Organizador: EventoOrganizadorDTO{
			ID: evento.Organizador.ID, UsuarioID: evento.Organizador.UsuarioID,
			RazaoSocial: evento.Organizador.RazaoSocial, NomeFantasia: evento.Organizador.NomeFantasia,
			CNPJ: evento.Organizador.CNPJ,
		},
		Logradouro: evento.Logradouro, Cidade: evento.Cidade, Estado: evento.Estado,
		CEP: evento.CEP, EventoLogo: evento.EventoLogo, Ingressos: ingressos,
	}
}
