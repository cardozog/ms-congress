package repositories

import (
	models "ms-congress/internal/models/usuario"

	"gorm.io/gorm"
)

type UsuarioRepository struct {
	db *gorm.DB
}

type UsuarioRepositoryInterface interface {
	BuscarUsuario(usuarioID uint) (*models.Usuario, error)
	BuscarPessoaFisica(usuarioID uint) (*models.PessoaFisica, error)
	BuscarPessoaJuridica(usuarioID uint) (*models.PessoaJuridica, error)
	BuscarImprensa(usuarioID uint) (*models.Imprensa, error)
	AtualizarSenha(usuarioID uint, senha string) error
	AtualizarPessoaFisica(usuarioID uint, dados map[string]interface{}) (*models.PessoaFisica, error)
	AtualizarPessoaJuridica(usuarioID uint, dados map[string]interface{}) (*models.PessoaJuridica, error)
	AtualizarImprensa(usuarioID uint, dados map[string]interface{}) (*models.Imprensa, error)
}

func NewUsuarioRepository(db *gorm.DB) UsuarioRepositoryInterface {
	return &UsuarioRepository{db: db}
}

func (r *UsuarioRepository) BuscarUsuario(usuarioID uint) (*models.Usuario, error) {
	var usuario models.Usuario
	err := r.db.
		Preload("TipoPessoa").
		Preload("Papeis.Role").
		Preload("Telefones").
		Preload("Enderecos").
		First(&usuario, usuarioID).Error
	if err != nil {
		return nil, err
	}
	return &usuario, nil
}

func (r *UsuarioRepository) BuscarPessoaFisica(usuarioID uint) (*models.PessoaFisica, error) {
	var pessoa models.PessoaFisica
	err := r.db.Where("usuario_id = ?", usuarioID).First(&pessoa).Error
	if err != nil {
		return nil, err
	}
	return &pessoa, nil
}

func (r *UsuarioRepository) BuscarPessoaJuridica(usuarioID uint) (*models.PessoaJuridica, error) {
	var pessoa models.PessoaJuridica
	err := r.db.Where("usuario_id = ?", usuarioID).First(&pessoa).Error
	if err != nil {
		return nil, err
	}
	return &pessoa, nil
}

func (r *UsuarioRepository) BuscarImprensa(usuarioID uint) (*models.Imprensa, error) {
	var imprensa models.Imprensa
	err := r.db.Where("usuario_id = ?", usuarioID).First(&imprensa).Error
	if err != nil {
		return nil, err
	}
	return &imprensa, nil
}

func (r *UsuarioRepository) AtualizarSenha(usuarioID uint, senha string) error {
	return r.db.Model(&models.Usuario{}).Where("id = ?", usuarioID).Update("senha", senha).Error
}

func (r *UsuarioRepository) AtualizarPessoaFisica(usuarioID uint, dados map[string]interface{}) (*models.PessoaFisica, error) {
	var pessoa models.PessoaFisica
	if err := r.db.Where("usuario_id = ?", usuarioID).First(&pessoa).Error; err != nil {
		return nil, err
	}
	if err := r.db.Model(&pessoa).Updates(dados).Error; err != nil {
		return nil, err
	}
	if err := r.db.First(&pessoa, pessoa.ID).Error; err != nil {
		return nil, err
	}
	return &pessoa, nil
}

func (r *UsuarioRepository) AtualizarPessoaJuridica(usuarioID uint, dados map[string]interface{}) (*models.PessoaJuridica, error) {
	var pessoa models.PessoaJuridica
	if err := r.db.Where("usuario_id = ?", usuarioID).First(&pessoa).Error; err != nil {
		return nil, err
	}
	if err := r.db.Model(&pessoa).Updates(dados).Error; err != nil {
		return nil, err
	}
	if err := r.db.First(&pessoa, pessoa.ID).Error; err != nil {
		return nil, err
	}
	return &pessoa, nil
}

func (r *UsuarioRepository) AtualizarImprensa(usuarioID uint, dados map[string]interface{}) (*models.Imprensa, error) {
	var imprensa models.Imprensa
	if err := r.db.Where("usuario_id = ?", usuarioID).First(&imprensa).Error; err != nil {
		return nil, err
	}
	if err := r.db.Model(&imprensa).Updates(dados).Error; err != nil {
		return nil, err
	}
	if err := r.db.First(&imprensa, imprensa.ID).Error; err != nil {
		return nil, err
	}
	return &imprensa, nil
}
