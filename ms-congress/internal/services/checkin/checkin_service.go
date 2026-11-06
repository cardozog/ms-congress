package services

import (
	"net/http"

	"ms-congress/internal/erros"
	repositories "ms-congress/internal/repositories/checkin"
)

type Service struct {
	repository *repositories.CheckinRepository
}

func NewCheckinService(repository *repositories.CheckinRepository) *Service {
	return &Service{repository: repository}
}

func (s *Service) Listar(eventoID uint64, organizadorUsuarioID uint, search string) ([]repositories.CheckinRecord, *erros.ApiError) {
	owned, err := s.repository.EventoPertenceAoOrganizador(eventoID, organizadorUsuarioID)
	if err != nil {
		return nil, &erros.ApiError{StatusCode: http.StatusInternalServerError, Error: err}
	}
	if !owned {
		return nil, &erros.ApiError{StatusCode: http.StatusNotFound, Error: erros.ErroEventoNaoEncontrado}
	}
	records, err := s.repository.Listar(eventoID, organizadorUsuarioID, search)
	if err != nil {
		return nil, &erros.ApiError{StatusCode: http.StatusInternalServerError, Error: err}
	}
	return records, nil
}
