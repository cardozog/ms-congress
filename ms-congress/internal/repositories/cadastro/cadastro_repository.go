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
	CadastrarCongressista(usuario *models.Usuario, pessoa *models.PessoaFisica) error
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

func (r *CadastroRepository) CadastrarCongressista(usuario *models.Usuario, pessoa *models.PessoaFisica) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var tipo models.TipoPessoa
		if err := tx.Where("tipo = ?", "FISICA").First(&tipo).Error; err != nil {
			return err
		}
		usuario.TipoPessoaID = tipo.ID
		if err := tx.Create(usuario).Error; err != nil {
			return err
		}
		pessoa.UsuarioID = usuario.ID
		if err := tx.Create(pessoa).Error; err != nil {
			return err
		}
		var role models.Role
		if err := tx.Where("nome = ?", "PARTICIPANTE").First(&role).Error; err != nil {
			return err
		}
		return tx.Create(&models.UsuarioRole{UsuarioID: usuario.ID, RoleID: role.ID}).Error
	})
}
