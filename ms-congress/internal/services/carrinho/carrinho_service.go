package services

import (
	"errors"
	"net/http"
	"time"

	"ms-congress/internal/erros"
	eventomodels "ms-congress/internal/models/evento"
	repositories "ms-congress/internal/repositories/carrinho"

	"gorm.io/gorm"
)

type CartEntry struct {
	EventoIngressoID uint64 `json:"evento_ingresso_id"`
	Quantidade       int    `json:"quantidade"`
}

type CartItem struct {
	ID            uint64
	Evento        eventomodels.Evento
	Ingresso      eventomodels.EventoIngresso
	Quantidade    int
	PrecoUnitario float64
}

type CartSnapshot struct {
	Itens []CartItem
}

type CarrinhoServiceInterface interface {
	Buscar([]CartEntry) (*CartSnapshot, *erros.ApiError)
	Adicionar([]CartEntry, uint64, int) ([]CartEntry, *erros.ApiError)
	Atualizar([]CartEntry, uint64, int) ([]CartEntry, *erros.ApiError)
	Remover([]CartEntry, uint64) ([]CartEntry, *erros.ApiError)
	Validar([]CartEntry) (*CartSnapshot, *erros.ApiError)
}

type CarrinhoService struct {
	repository repositories.CarrinhoRepositoryInterface
}

func NewCarrinhoService(repository repositories.CarrinhoRepositoryInterface) CarrinhoServiceInterface {
	return &CarrinhoService{repository: repository}
}

func (s *CarrinhoService) Buscar(entries []CartEntry) (*CartSnapshot, *erros.ApiError) {
	return s.Validar(entries)
}

func (s *CarrinhoService) Validar(entries []CartEntry) (*CartSnapshot, *erros.ApiError) {
	snapshot := &CartSnapshot{Itens: make([]CartItem, 0, len(entries))}
	for _, entry := range entries {
		lote, apiErr := s.validarLote(entry.EventoIngressoID)
		if apiErr != nil {
			return nil, apiErr
		}
		if entry.Quantidade <= 0 || entry.Quantidade > disponivel(lote) {
			return nil, cartValidation("quantidade solicitada indisponível")
		}
		snapshot.Itens = append(snapshot.Itens, CartItem{
			ID: entry.EventoIngressoID, Evento: lote.Evento, Ingresso: *lote,
			Quantidade: entry.Quantidade, PrecoUnitario: lote.Preco,
		})
	}
	return snapshot, nil
}

func (s *CarrinhoService) Adicionar(entries []CartEntry, loteID uint64, quantidade int) ([]CartEntry, *erros.ApiError) {
	if quantidade <= 0 {
		return nil, cartValidation("quantidade deve ser maior que zero")
	}
	lote, apiErr := s.validarLote(loteID)
	if apiErr != nil {
		return nil, apiErr
	}
	for index := range entries {
		if entries[index].EventoIngressoID == loteID {
			entries[index].Quantidade += quantidade
			if entries[index].Quantidade > disponivel(lote) {
				return nil, cartValidation("quantidade solicitada indisponível")
			}
			return entries, nil
		}
	}
	if quantidade > disponivel(lote) {
		return nil, cartValidation("quantidade solicitada indisponível")
	}
	return append(entries, CartEntry{EventoIngressoID: loteID, Quantidade: quantidade}), nil
}

func (s *CarrinhoService) Atualizar(entries []CartEntry, itemID uint64, quantidade int) ([]CartEntry, *erros.ApiError) {
	if quantidade <= 0 {
		return nil, cartValidation("quantidade deve ser maior que zero")
	}
	lote, apiErr := s.validarLote(itemID)
	if apiErr != nil {
		return nil, apiErr
	}
	for index := range entries {
		if entries[index].EventoIngressoID == itemID {
			if quantidade > disponivel(lote) {
				return nil, cartValidation("quantidade solicitada indisponível")
			}
			entries[index].Quantidade = quantidade
			return entries, nil
		}
	}
	return nil, cartNotFound()
}

func (s *CarrinhoService) Remover(entries []CartEntry, itemID uint64) ([]CartEntry, *erros.ApiError) {
	for index := range entries {
		if entries[index].EventoIngressoID == itemID {
			return append(entries[:index], entries[index+1:]...), nil
		}
	}
	return nil, cartNotFound()
}

func (s *CarrinhoService) validarLote(id uint64) (*eventomodels.EventoIngresso, *erros.ApiError) {
	lote, err := s.repository.BuscarLote(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, cartNotFound()
	}
	if err != nil {
		return nil, internalCartError(err)
	}
	now := time.Now()
	if !lote.Evento.Publicado || lote.Evento.DataFim.Before(now) ||
		lote.DataInicioVenda.After(now) || lote.DataFimVenda.Before(now) {
		return nil, cartValidation("ingresso fora do período de venda")
	}
	if disponivel(lote) <= 0 {
		return nil, cartValidation("ingresso esgotado")
	}
	return lote, nil
}

func disponivel(lote *eventomodels.EventoIngresso) int {
	if len(lote.Ingressos) == 0 {
		return lote.Quantidade
	}
	total := 0
	for _, ingresso := range lote.Ingressos {
		if ingresso.Status == "DISPONIVEL" {
			total++
		}
	}
	return total
}

func cartValidation(message string) *erros.ApiError {
	return &erros.ApiError{StatusCode: http.StatusBadRequest, Error: errors.New(message)}
}
func cartNotFound() *erros.ApiError {
	return &erros.ApiError{StatusCode: http.StatusNotFound, Error: errors.New("item do carrinho não encontrado")}
}
func internalCartError(err error) *erros.ApiError {
	return &erros.ApiError{StatusCode: http.StatusInternalServerError, Error: err}
}
