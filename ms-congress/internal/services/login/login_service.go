package services

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"ms-congress/internal/erros"
	repositories "ms-congress/internal/repositories/login"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginServiceInterface interface {
	Login(email, senha, audience string) (*LoginResponse, *erros.ApiError)
}

type LoginService struct {
	repository repositories.LoginRepositoryInterface
	jwtSecret  string
}

type LoginResponse struct {
	Token        string   `json:"token"`
	TokenType    string   `json:"tokenType"`
	ExpiresIn    int64    `json:"expiresIn"`
	ID           uint     `json:"id"`
	Email        string   `json:"email"`
	TipoPessoaID uint     `json:"tipoPessoaID"`
	Papeis       []string `json:"papeis"`
}

func NewLoginService(repository repositories.LoginRepositoryInterface, jwtSecret string) LoginServiceInterface {
	return &LoginService{repository: repository, jwtSecret: jwtSecret}
}

func (s *LoginService) Login(email, senha, audience string) (*LoginResponse, *erros.ApiError) {
	email = strings.TrimSpace(email)
	if email == "" || senha == "" {
		return nil, &erros.ApiError{StatusCode: http.StatusBadRequest, Error: erros.ErroCredenciaisObrigatorias}
	}

	usuario, err := s.repository.BuscarPorEmail(email)
	if errors.Is(err, gorm.ErrRecordNotFound) || (err == nil && bcrypt.CompareHashAndPassword([]byte(usuario.Senha), []byte(senha)) != nil) {
		return nil, &erros.ApiError{StatusCode: http.StatusUnauthorized, Error: erros.ErroCredenciaisInvalidas}
	}
	if err != nil {
		return nil, &erros.ApiError{StatusCode: http.StatusInternalServerError, Error: err}
	}
	papeis := make([]string, 0, len(usuario.Papeis))
	for _, papel := range usuario.Papeis {
		papeis = append(papeis, papel.Role.Nome)
	}
	if audience == "congressista" && !containsRole(papeis, "PARTICIPANTE") {
		return nil, &erros.ApiError{StatusCode: http.StatusUnauthorized, Error: erros.ErroCredenciaisInvalidas}
	}
	if audience == "admin" && containsRole(papeis, "PARTICIPANTE") && !containsRole(papeis, "ADMIN") && !containsRole(papeis, "ORGANIZADOR") && !containsRole(papeis, "IMPRENSA") {
		return nil, &erros.ApiError{StatusCode: http.StatusUnauthorized, Error: erros.ErroCredenciaisInvalidas}
	}

	if strings.TrimSpace(s.jwtSecret) == "" {
		return nil, &erros.ApiError{StatusCode: http.StatusInternalServerError, Error: erros.ErroConfiguracaoToken}
	}

	agora := time.Now()
	expiraEm := agora.Add(time.Hour)
	claims := jwt.MapClaims{
		"sub":          strconv.FormatUint(uint64(usuario.ID), 10),
		"email":        usuario.Email,
		"tipoPessoaID": usuario.TipoPessoaID,
		"papeis":       papeis,
		"iat":          agora.Unix(),
		"exp":          expiraEm.Unix(),
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, &erros.ApiError{StatusCode: http.StatusInternalServerError, Error: err}
	}

	return &LoginResponse{
		Token: token, TokenType: "Bearer", ExpiresIn: int64(time.Hour.Seconds()),
		ID: usuario.ID, Email: usuario.Email, TipoPessoaID: usuario.TipoPessoaID,
		Papeis: papeis,
	}, nil
}

func containsRole(roles []string, expected string) bool {
	for _, role := range roles {
		if role == expected {
			return true
		}
	}
	return false
}
