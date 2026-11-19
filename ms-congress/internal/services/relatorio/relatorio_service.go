package services

import (
	"net/http"
	"time"

	"ms-congress/internal/erros"
	repositories "ms-congress/internal/repositories/relatorio"
)

type Service struct {
	repository *repositories.RelatorioRepository
}

func NewService(repository *repositories.RelatorioRepository) *Service {
	return &Service{repository: repository}
}

func (s *Service) ListarEventos(organizadorUsuarioID uint, organizadorID uint64, dataDe, dataAte *time.Time) ([]repositories.EventoRelatorio, *erros.ApiError) {
	relatorios, err := s.repository.ListarEventos(organizadorUsuarioID, organizadorID, dataDe, dataAte)
	if err != nil {
		return nil, &erros.ApiError{StatusCode: http.StatusInternalServerError, Error: err}
	}
	return relatorios, nil
}
