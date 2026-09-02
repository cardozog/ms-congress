package db

import (
	"context"
	"log"
	models "ms-congress/internal/models/usuario"
	"time"

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
	log.Default().Println("Migrations applied successfully")
	return database.Where(Migration{Version: 1}).FirstOrCreate(&Migration{
		Version: 1,
		Name:    "initialize_schema_migrations",
	}).Error
}
