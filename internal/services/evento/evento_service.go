package services

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"ms-congress/internal/erros"
	models "ms-congress/internal/models/evento"
	repositories "ms-congress/internal/repositories/evento"

	"gorm.io/gorm"
)

type EventoServiceInterface interface {
	Criar(evento *models.Evento) *erros.ApiError
	BuscarPorID(id uint64) (*models.Evento, *erros.ApiError)
	BuscarPorOrganizador(organizadorID uint64) ([]models.Evento, *erros.ApiError)
	Atualizar(evento *models.Evento) *erros.ApiError
	Excluir(id uint64) *erros.ApiError
}

type EventoService struct {
	repository repositories.EventoRepositoryInterface
}

func NewEventoService(repository repositories.EventoRepositoryInterface) EventoServiceInterface {
	return &EventoService{repository: repository}
}

func (s *EventoService) Criar(evento *models.Evento) *erros.ApiError {
	if apiErr := validarEvento(evento, false); apiErr != nil {
		return apiErr
	}
	if err := s.repository.Create(evento); err != nil {
		return internalError(err)
	}
	return nil
}

func (s *EventoService) BuscarPorID(id uint64) (*models.Evento, *erros.ApiError) {
	evento, err := s.repository.GetByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, eventoNotFound()
	}
	if err != nil {
		return nil, internalError(err)
	}
	return evento, nil
}

func (s *EventoService) BuscarPorOrganizador(organizadorID uint64) ([]models.Evento, *erros.ApiError) {
	eventos, err := s.repository.BuscarEventosPorOrganizador(organizadorID)
	if err != nil {
		return nil, internalError(err)
	}
	return eventos, nil
}

func (s *EventoService) Atualizar(evento *models.Evento) *erros.ApiError {
	if apiErr := validarEvento(evento, true); apiErr != nil {
		return apiErr
	}
	if _, apiErr := s.BuscarPorID(evento.ID); apiErr != nil {
		return apiErr
	}
	if err := s.repository.Update(evento); err != nil {
		return internalError(err)
	}
	return nil
}

func (s *EventoService) Excluir(id uint64) *erros.ApiError {
	if id == 0 {
		return eventoValidationError("id é obrigatório")
	}
	evento, apiErr := s.BuscarPorID(id)
	if apiErr != nil {
		return apiErr
	}
	if err := s.repository.Delete(evento); err != nil {
		return internalError(err)
	}
	return nil
}

func eventoNotFound() *erros.ApiError {
	return &erros.ApiError{StatusCode: http.StatusNotFound, Error: erros.ErroEventoNaoEncontrado}
}

func internalError(err error) *erros.ApiError {
	return &erros.ApiError{StatusCode: http.StatusInternalServerError, Error: err}
}

func validarEvento(evento *models.Evento, validarID bool) *erros.ApiError {
	if evento == nil {
		return eventoValidationError("evento é obrigatório")
	}
	if validarID && evento.ID == 0 {
		return eventoValidationError("id é obrigatório")
	}
	if strings.TrimSpace(evento.Nome) == "" {
		return eventoValidationError("nome é obrigatório")
	}
	if strings.TrimSpace(evento.Descricao) == "" {
		return eventoValidationError("descricao é obrigatório")
	}
	if evento.DataInicio.IsZero() {
		return eventoValidationError("dataInicio é obrigatória")
	}
	if evento.DataFim.IsZero() {
		return eventoValidationError("dataFim é obrigatória")
	}
	if !evento.DataFim.After(evento.DataInicio) {
		return eventoValidationError("dataFim deve ser posterior a dataInicio")
	}
	if evento.OrganizadorID == 0 {
		return eventoValidationError("organizadorID é obrigatório")
	}
	if strings.TrimSpace(evento.Logradouro) == "" {
		return eventoValidationError("logradouro é obrigatório")
	}
	if strings.TrimSpace(evento.Cidade) == "" {
		return eventoValidationError("cidade é obrigatória")
	}
	if strings.TrimSpace(evento.Estado) == "" {
		return eventoValidationError("estado é obrigatório")
	}
	if strings.TrimSpace(evento.CEP) == "" {
		return eventoValidationError("cep é obrigatório")
	}
	return nil
}

func eventoValidationError(message string) *erros.ApiError {
	return &erros.ApiError{StatusCode: http.StatusBadRequest, Error: fmt.Errorf("evento inválido: %s", message)}
}
