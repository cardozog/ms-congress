package services

import (
	"errors"
	"net/http"
	"strings"

	"ms-congress/internal/erros"
	pedidomodels "ms-congress/internal/models/pedido"
	repositories "ms-congress/internal/repositories/pedido"
	carrinho "ms-congress/internal/services/carrinho"
)

type PedidoService struct {
	repository *repositories.PedidoRepository
}

func NewPedidoService(repository *repositories.PedidoRepository) *PedidoService {
	return &PedidoService{repository: repository}
}

func (s *PedidoService) BuscarPorChave(usuarioID uint, key string) (*pedidomodels.Pedido, error) {
	return s.repository.BuscarPorChave(usuarioID, strings.TrimSpace(key))
}

func (s *PedidoService) Criar(usuarioID uint, key string, entries []carrinho.CartEntry) (interface{}, *erros.ApiError) {
	if usuarioID == 0 || strings.TrimSpace(key) == "" {
		return nil, &erros.ApiError{StatusCode: http.StatusBadRequest, Error: errors.New("Idempotency-Key é obrigatório")}
	}
	if len(entries) == 0 {
		return nil, &erros.ApiError{StatusCode: http.StatusBadRequest, Error: errors.New("carrinho vazio")}
	}
	pedido, err := s.repository.Criar(usuarioID, strings.TrimSpace(key), entries)
	if err != nil {
		return nil, &erros.ApiError{StatusCode: http.StatusBadRequest, Error: err}
	}
	return pedido, nil
}
