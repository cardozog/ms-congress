package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	models "ms-congress/internal/models/usuario"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Migration struct {
	Version   uint      `gorm:"primaryKey"`
	Name      string    `gorm:"not null"`
	AppliedAt time.Time `gorm:"not null;autoCreateTime"`
}

func (Migration) TableName() string {
	return "schema_migrations"
}

func RunMigrations(ctx context.Context, database *gorm.DB) error {
	database = database.WithContext(ctx)
	if err := database.AutoMigrate(&Migration{}); err != nil {
		return err
	}

	if err := database.AutoMigrate(
		&models.Usuario{},
		&models.TipoPessoa{},
		&models.PessoaFisica{},
		&models.PessoaJuridica{},
		&models.Role{},
		&models.UsuarioRole{},
		&models.TelefoneUsuario{},
		&models.EnderecoUsuario{},
		&models.Imprensa{}); err != nil {
		return err
	}
	if err := seedAdmin(database); err != nil {
		return err
	}
	log.Default().Println("Migrations applied successfully")
	return database.Where(Migration{Version: 1}).FirstOrCreate(&Migration{
		Version: 1,
		Name:    "initialize_schema_migrations",
	}).Error
}

func seedAdmin(database *gorm.DB) error {
	email := os.Getenv("ADMIN_EMAIL")
	password := os.Getenv("ADMIN_PASSWORD")
	if email == "" || password == "" {
		return fmt.Errorf("ADMIN_EMAIL and ADMIN_PASSWORD must be set")
	}

	return database.Transaction(func(tx *gorm.DB) error {
		tipoPessoa := models.TipoPessoa{Tipo: "FISICA"}
		if err := tx.Where("tipo = ?", tipoPessoa.Tipo).FirstOrCreate(&tipoPessoa).Error; err != nil {
			return err
		}

		role := models.Role{Nome: "ADMIN"}
		if err := tx.Where("nome = ?", role.Nome).FirstOrCreate(&role).Error; err != nil {
			return err
		}

		usuario := models.Usuario{}
		err := tx.Where("email = ?", email).First(&usuario).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			senha, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				return err
			}
			usuario = models.Usuario{
				Email:        email,
				Senha:        string(senha),
				TipoPessoaID: tipoPessoa.ID,
			}
			if err := tx.Create(&usuario).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err
		}

		return tx.Where(models.UsuarioRole{
			UsuarioID: usuario.ID,
			RoleID:    role.ID,
		}).FirstOrCreate(&models.UsuarioRole{}).Error
	})
}
