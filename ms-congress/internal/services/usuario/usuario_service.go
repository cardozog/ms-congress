package services

import (
	"errors"
	"net/http"
	"strings"

	"ms-congress/internal/erros"
	models "ms-congress/internal/models/usuario"
	repositories "ms-congress/internal/repositories/usuario"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UsuarioServiceInterface interface {
	BuscarDetalhes(usuarioID uint) (*models.UsuarioDetalhadoResponseDTO, *erros.ApiError)
	AtualizarPessoaFisica(usuarioID uint, dados PessoaFisicaUpdate) (*models.PessoaFisica, *erros.ApiError)
	AtualizarPessoaJuridica(usuarioID uint, dados PessoaJuridicaUpdate) (*models.PessoaJuridica, *erros.ApiError)
	AtualizarImprensa(usuarioID uint, dados ImprensaUpdate) (*models.Imprensa, *erros.ApiError)
}

type UsuarioService struct {
	repository repositories.UsuarioRepositoryInterface
}

type PessoaFisicaUpdate struct {
	Nome           string
	DataNascimento string
	Senha          string
}

type PessoaJuridicaUpdate struct {
	RazaoSocial  string
	NomeFantasia string
	Senha        string
}

type ImprensaUpdate struct {
	NomeVeiculo string
	Site        string
	Senha       string
}

func NewUsuarioService(repository repositories.UsuarioRepositoryInterface) UsuarioServiceInterface {
	return &UsuarioService{repository: repository}
}

func (s *UsuarioService) BuscarDetalhes(usuarioID uint) (*models.UsuarioDetalhadoResponseDTO, *erros.ApiError) {
	if apiErr := validarUsuarioID(usuarioID); apiErr != nil {
		return nil, apiErr
	}

	usuario, err := s.repository.BuscarUsuario(usuarioID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, usuarioNaoEncontrado()
	}
	if err != nil {
		return nil, erroInterno(err)
	}

	resposta := &models.UsuarioDetalhadoResponseDTO{
		ID:         usuario.ID,
		Email:      usuario.Email,
		TipoPessoa: models.TipoPessoaResponseDTO{ID: usuario.TipoPessoa.ID, Tipo: usuario.TipoPessoa.Tipo},
		Papeis:     make([]models.RoleResponseDTO, 0, len(usuario.Papeis)),
		Telefones:  make([]models.TelefoneResponseDTO, 0, len(usuario.Telefones)),
		Enderecos:  make([]models.EnderecoResponseDTO, 0, len(usuario.Enderecos)),
	}
	for _, papel := range usuario.Papeis {
		resposta.Papeis = append(resposta.Papeis, models.RoleResponseDTO{ID: papel.Role.ID, Nome: papel.Role.Nome})
	}
	for _, telefone := range usuario.Telefones {
		resposta.Telefones = append(resposta.Telefones, models.TelefoneResponseDTO{ID: telefone.ID, Numero: telefone.Numero})
	}
	for _, endereco := range usuario.Enderecos {
		resposta.Enderecos = append(resposta.Enderecos, models.EnderecoResponseDTO{
			ID: endereco.ID, Logradouro: endereco.Logradouro, Cidade: endereco.Cidade,
			Estado: endereco.Estado, CEP: endereco.CEP, EnderecoCobranca: endereco.EnderecoCobranca,
		})
	}

	if pessoa, err := s.repository.BuscarPessoaFisica(usuarioID); err == nil {
		resposta.PessoaFisica = &models.PessoaFisicaResponseDTO{ID: pessoa.ID, UsuarioID: pessoa.UsuarioID, Nome: pessoa.Nome, DataNascimento: pessoa.DataNascimento}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, erroInterno(err)
	}
	if pessoa, err := s.repository.BuscarPessoaJuridica(usuarioID); err == nil {
		resposta.PessoaJuridica = &models.PessoaJuridicaResponseDTO{ID: pessoa.ID, UsuarioID: pessoa.UsuarioID, RazaoSocial: pessoa.RazaoSocial, NomeFantasia: pessoa.NomeFantasia}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, erroInterno(err)
	}
	if imprensa, err := s.repository.BuscarImprensa(usuarioID); err == nil {
		resposta.Imprensa = &models.ImprensaResponseDTO{ID: imprensa.ID, UsuarioID: imprensa.UsuarioID, NomeVeiculo: imprensa.NomeVeiculo, Site: imprensa.Site}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, erroInterno(err)
	}

	return resposta, nil
}

func (s *UsuarioService) AtualizarPessoaFisica(usuarioID uint, dados PessoaFisicaUpdate) (*models.PessoaFisica, *erros.ApiError) {
	if apiErr := validarUsuarioID(usuarioID); apiErr != nil {
		return nil, apiErr
	}
	if strings.TrimSpace(dados.Nome) == "" {
		return nil, erroValidacao("nome é obrigatório")
	}
	updates := map[string]interface{}{"nome": dados.Nome, "data_nascimento": dados.DataNascimento}
	if dados.Senha != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(dados.Senha), bcrypt.DefaultCost)
		if err != nil {
			return nil, erroInterno(err)
		}
		if err := s.repository.AtualizarSenha(usuarioID, string(hash)); err != nil {
			return nil, erroInterno(err)
		}
	}
	pessoa, err := s.repository.AtualizarPessoaFisica(usuarioID, updates)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, usuarioNaoEncontrado()
	}
	if err != nil {
		return nil, erroInterno(err)
	}
	return pessoa, nil
}

func (s *UsuarioService) AtualizarPessoaJuridica(usuarioID uint, dados PessoaJuridicaUpdate) (*models.PessoaJuridica, *erros.ApiError) {
	if apiErr := validarUsuarioID(usuarioID); apiErr != nil {
		return nil, apiErr
	}
	if strings.TrimSpace(dados.RazaoSocial) == "" {
		return nil, erroValidacao("razaoSocial é obrigatória")
	}
	updates := map[string]interface{}{"razao_social": dados.RazaoSocial, "nome_fantasia": dados.NomeFantasia}
	if dados.Senha != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(dados.Senha), bcrypt.DefaultCost)
		if err != nil {
			return nil, erroInterno(err)
		}
		if err := s.repository.AtualizarSenha(usuarioID, string(hash)); err != nil {
			return nil, erroInterno(err)
		}
	}
	pessoa, err := s.repository.AtualizarPessoaJuridica(usuarioID, updates)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, usuarioNaoEncontrado()
	}
	if err != nil {
		return nil, erroInterno(err)
	}
	return pessoa, nil
}

func (s *UsuarioService) AtualizarImprensa(usuarioID uint, dados ImprensaUpdate) (*models.Imprensa, *erros.ApiError) {
	if apiErr := validarUsuarioID(usuarioID); apiErr != nil {
		return nil, apiErr
	}
	if strings.TrimSpace(dados.NomeVeiculo) == "" && strings.TrimSpace(dados.Site) == "" && dados.Senha == "" {
		return nil, erroValidacao("informe ao menos um dado para atualização")
	}
	updates := map[string]interface{}{"nome_veiculo": dados.NomeVeiculo, "site": dados.Site}
	if dados.Senha != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(dados.Senha), bcrypt.DefaultCost)
		if err != nil {
			return nil, erroInterno(err)
		}
		if err := s.repository.AtualizarSenha(usuarioID, string(hash)); err != nil {
			return nil, erroInterno(err)
		}
	}
	imprensa, err := s.repository.AtualizarImprensa(usuarioID, updates)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, usuarioNaoEncontrado()
	}
	if err != nil {
		return nil, erroInterno(err)
	}
	return imprensa, nil
}

func validarUsuarioID(usuarioID uint) *erros.ApiError {
	if usuarioID == 0 {
		return &erros.ApiError{StatusCode: http.StatusBadRequest, Error: errors.New("usuarioID é obrigatório")}
	}
	return nil
}

func erroValidacao(message string) *erros.ApiError {
	return &erros.ApiError{StatusCode: http.StatusBadRequest, Error: errors.New(message)}
}

func usuarioNaoEncontrado() *erros.ApiError {
	return &erros.ApiError{StatusCode: http.StatusNotFound, Error: errors.New("dados do usuário não encontrados")}
}

func erroInterno(err error) *erros.ApiError {
	return &erros.ApiError{StatusCode: http.StatusInternalServerError, Error: err}
}
