package services

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

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
	Publicar(id uint64) *erros.ApiError
	Encerrar(id uint64) *erros.ApiError
	Excluir(id uint64) *erros.ApiError
	ListarTiposIngresso() ([]models.TipoIngresso, *erros.ApiError)
}

func (s *EventoService) ListarTiposIngresso() ([]models.TipoIngresso, *erros.ApiError) {
	tipos, err := s.repository.ListarTiposIngresso()
	if err != nil {
		return nil, internalError(err)
	}
	return tipos, nil
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
	eventoAtual, apiErr := s.BuscarPorID(evento.ID)
	if apiErr != nil {
		return apiErr
	}
	evento.Publicado = eventoAtual.Publicado
	if err := s.repository.Update(evento); err != nil {
		return internalError(err)
	}
	return nil
}

func (s *EventoService) Publicar(id uint64) *erros.ApiError {
	evento, apiErr := s.BuscarPorID(id)
	if apiErr != nil {
		return apiErr
	}
	if evento.Publicado {
		return nil
	}
	if err := s.repository.Publicar(id); err != nil {
		return internalError(err)
	}
	return nil
}

func (s *EventoService) Encerrar(id uint64) *erros.ApiError {
	evento, apiErr := s.BuscarPorID(id)
	if apiErr != nil {
		return apiErr
	}
	if !evento.Publicado {
		return &erros.ApiError{StatusCode: http.StatusConflict, Error: erros.ErroEventoNaoPublicado}
	}
	if evento.DataFim.Before(time.Now()) {
		return nil
	}
	if err := s.repository.Encerrar(id, time.Now()); err != nil {
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
	if evento.Publicado {
		return &erros.ApiError{StatusCode: http.StatusConflict, Error: erros.ErroEventoPublicado}
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
	if len(evento.Ingressos) == 0 {
		return eventoValidationError("ao menos um ingresso é obrigatório")
	}
	tipos := make(map[uint64]struct{}, len(evento.Ingressos))
	for _, ingresso := range evento.Ingressos {
		if ingresso.TipoIngressoID == 0 {
			return eventoValidationError("tipo_ingresso_id é obrigatório")
		}
		if _, existe := tipos[ingresso.TipoIngressoID]; existe {
			return eventoValidationError("não é permitido repetir o tipo de ingresso")
		}
		tipos[ingresso.TipoIngressoID] = struct{}{}
		if ingresso.Preco < 0 {
			return eventoValidationError("preco não pode ser negativo")
		}
		if ingresso.Quantidade <= 0 {
			return eventoValidationError("quantidade deve ser maior que zero")
		}
		if !ingresso.DataFimVenda.After(ingresso.DataInicioVenda) {
			return eventoValidationError("data_fim_venda deve ser posterior a data_inicio_venda")
		}
	}
	return nil
}

func eventoValidationError(message string) *erros.ApiError {
	return &erros.ApiError{StatusCode: http.StatusBadRequest, Error: fmt.Errorf("evento inválido: %s", message)}
}
