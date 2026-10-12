package services

import (
	"testing"

	models "ms-congress/internal/models/usuario"
)

type cadastroRepositoryFake struct {
	emailExists   bool
	createdUser   *models.Usuario
	createdPerson *models.PessoaFisica
}

func (f *cadastroRepositoryFake) EmailExiste(string) (bool, error)                     { return f.emailExists, nil }
func (f *cadastroRepositoryFake) CadastrarPessoaJuridica(*models.PessoaJuridica) error { return nil }
func (f *cadastroRepositoryFake) CadastrarCongressista(user *models.Usuario, person *models.PessoaFisica) error {
	f.createdUser, f.createdPerson = user, person
	return nil
}

func TestCadastrarCongressistaHashesPasswordAndAssignsData(t *testing.T) {
	repository := &cadastroRepositoryFake{}
	service := NewCadastroService(repository)

	if apiErr := service.CadastrarCongressista(" pessoa@example.com ", "password123", "Pessoa Teste", "123", "2000-01-01"); apiErr != nil {
		t.Fatalf("expected successful registration, got %v", apiErr)
	}
	if repository.createdUser == nil || repository.createdPerson == nil {
		t.Fatal("expected registration repository to receive user and person")
	}
	if repository.createdUser.Senha == "password123" || repository.createdUser.Senha == "" {
		t.Fatal("expected password to be hashed")
	}
	if repository.createdPerson.Nome != "Pessoa Teste" {
		t.Fatalf("unexpected person name: %q", repository.createdPerson.Nome)
	}
}

func TestCadastrarCongressistaRejectsExistingEmail(t *testing.T) {
	service := NewCadastroService(&cadastroRepositoryFake{emailExists: true})

	if apiErr := service.CadastrarCongressista("pessoa@example.com", "password123", "Pessoa", "123", ""); apiErr == nil || apiErr.StatusCode != 409 {
		t.Fatalf("expected conflict for existing email, got %#v", apiErr)
	}
}
