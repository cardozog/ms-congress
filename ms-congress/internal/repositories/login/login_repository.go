package repositories

import (
	models "ms-congress/internal/models/usuario"

	"gorm.io/gorm"
)

type LoginRepository struct {
	db *gorm.DB
}

type LoginRepositoryInterface interface {
	BuscarPorEmail(email string) (*models.Usuario, error)
}

func NewLoginRepository(db *gorm.DB) LoginRepositoryInterface {
	return &LoginRepository{db: db}
}

func (r *LoginRepository) BuscarPorEmail(email string) (*models.Usuario, error) {
	var usuario models.Usuario
	if err := r.db.Where("email = ?", email).First(&usuario).Error; err != nil {
		return nil, err
	}
	return &usuario, nil
}
