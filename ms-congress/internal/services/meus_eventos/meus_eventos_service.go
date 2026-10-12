package services

import (
	"net/http"

	"ms-congress/internal/erros"
	eventomodels "ms-congress/internal/models/evento"
	repositories "ms-congress/internal/repositories/meus_eventos"
)

type Service struct {
	repository *repositories.MeusEventosRepository
}

func NewService(repository *repositories.MeusEventosRepository) *Service {
	return &Service{repository: repository}
}

func (s *Service) Listar(usuarioID uint) ([]eventomodels.Ingresso, *erros.ApiError) {
	ingressos, err := s.repository.Listar(usuarioID)
	if err != nil {
		return nil, &erros.ApiError{StatusCode: http.StatusInternalServerError, Error: err}
	}
	return ingressos, nil
}

func (s *Service) ListarPorEvento(usuarioID uint, eventoID uint64) ([]eventomodels.Ingresso, *erros.ApiError) {
	ingressos, err := s.repository.ListarPorEvento(usuarioID, eventoID)
	if err != nil {
		return nil, &erros.ApiError{StatusCode: http.StatusInternalServerError, Error: err}
	}
	return ingressos, nil
}
