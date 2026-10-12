package services

import (
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
	"ms-congress/internal/erros"
	models "ms-congress/internal/models/usuario"
	cadastroRepository "ms-congress/internal/repositories/cadastro"
	"net/http"
	"strings"
	"unicode/utf8"
)

type CadastroServiceInterface interface {
	EmailExiste(email string) *erros.ApiError
	CadastrarCongressista(email, senha, nome, cpf, dataNascimento string) *erros.ApiError
}

func (s *CadastroService) CadastrarCongressista(email, senha, nome, cpf, dataNascimento string) *erros.ApiError {
	email, nome, cpf = strings.TrimSpace(email), strings.TrimSpace(nome), strings.TrimSpace(cpf)
	if email == "" || senha == "" || nome == "" || cpf == "" || !strings.Contains(email, "@") || utf8.RuneCountInString(senha) < 8 {
		return &erros.ApiError{StatusCode: http.StatusBadRequest, Error: erros.ErroDadosCadastroInvalidos}
	}
	if err := s.EmailExiste(email); err != nil {
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(senha), bcrypt.DefaultCost)
	if err != nil {
		return &erros.ApiError{StatusCode: http.StatusInternalServerError, Error: err}
	}
	err = s.repository.CadastrarCongressista(
		&models.Usuario{Email: email, Senha: string(hash)},
		&models.PessoaFisica{Nome: nome, CPF: cpf, DataNascimento: strings.TrimSpace(dataNascimento)},
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return &erros.ApiError{StatusCode: http.StatusConflict, Error: erros.ErroDadosJaCadastrados}
		}
		return &erros.ApiError{StatusCode: http.StatusInternalServerError, Error: err}
	}
	return nil
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
