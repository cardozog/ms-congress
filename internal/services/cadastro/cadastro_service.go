package services

import (
	"ms-congress/internal/erros"
	cadastroRepository "ms-congress/internal/repositories/cadastro"
	"net/http"
)

type CadastroServiceInterface interface {
	EmailExiste(email string) *erros.ApiError
}

type CadastroService struct {
	repository cadastroRepository.CadastroRepositoryInterface
}

func NewCadastroService(repository cadastroRepository.CadastroRepositoryInterface) CadastroServiceInterface {
	return &CadastroService{
		repository: repository,
	}
}

func (s *CadastroService) EmailExiste(email string) *erros.ApiError {

	existe, err := s.repository.EmailExiste(email)
	if err != nil {
		return &erros.ApiError{
			StatusCode: 500,
			Error:      err,
		}
	}

	if existe {
		return &erros.ApiError{
			StatusCode: http.StatusConflict,
			Error:      erros.ErroEmailJaCadastrado,
		}
	}
	return nil
}
