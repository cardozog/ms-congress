package repositories

import (
	models "ms-congress/internal/models/usuario"

	"gorm.io/gorm"
)

type CadastroRepository struct {
	db *gorm.DB
}

type CadastroRepositoryInterface interface {
	EmailExiste(email string) (bool, error)
	CadastrarPessoaJuridica(pj *models.PessoaJuridica) error
}

func NewCadastroRepository(db *gorm.DB) CadastroRepositoryInterface {
	return &CadastroRepository{
		db: db,
	}
}

func (r *CadastroRepository) EmailExiste(email string) (bool, error) {

	var count int64
	if err := r.db.Model(&models.Usuario{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil

}

func (r *CadastroRepository) CadastrarPessoaJuridica(pj *models.PessoaJuridica) error {
	if err := r.db.Create(pj).Error; err != nil {
		return err
	}
	return nil
}
