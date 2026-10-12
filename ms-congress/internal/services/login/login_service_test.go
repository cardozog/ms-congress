package services

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
	models "ms-congress/internal/models/usuario"
)

type loginRepositoryFake struct {
	user *models.Usuario
}

func hashPasswordForTest(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(hash), err
}

func (f *loginRepositoryFake) BuscarPorEmail(string) (*models.Usuario, error) { return f.user, nil }

func TestLoginCongressistaRequiresParticipantRole(t *testing.T) {
	hash, err := hashPasswordForTest("password123")
	if err != nil {
		t.Fatal(err)
	}
	repository := &loginRepositoryFake{user: &models.Usuario{ID: 1, Email: "pessoa@example.com", Senha: hash, Papeis: []models.UsuarioRole{{Role: models.Role{Nome: "PARTICIPANTE"}}}}}
	service := NewLoginService(repository, "secret")

	response, apiErr := service.Login("pessoa@example.com", "password123", "congressista")
	if apiErr != nil || response == nil || len(response.Papeis) != 1 {
		t.Fatalf("expected congressista login, response=%#v error=%#v", response, apiErr)
	}
}

func TestLoginCongressistaRejectsAdministrativeOnlyUser(t *testing.T) {
	hash, err := hashPasswordForTest("password123")
	if err != nil {
		t.Fatal(err)
	}
	service := NewLoginService(&loginRepositoryFake{user: &models.Usuario{ID: 1, Email: "admin@example.com", Senha: hash, Papeis: []models.UsuarioRole{{Role: models.Role{Nome: "ADMIN"}}}}}, "secret")

	if _, apiErr := service.Login("admin@example.com", "password123", "congressista"); apiErr == nil || apiErr.StatusCode != 401 {
		t.Fatalf("expected unauthorized congressista login, got %#v", apiErr)
	}
}
