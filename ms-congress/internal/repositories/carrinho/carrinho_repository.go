package repositories

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	carrinhomodels "ms-congress/internal/models/carrinho"
	eventomodels "ms-congress/internal/models/evento"
)

type CarrinhoRepositoryInterface interface {
	BuscarOuCriar(usuarioID uint) (*carrinhomodels.Carrinho, error)
	Buscar(usuarioID uint) (*carrinhomodels.Carrinho, error)
	BuscarLote(id uint64) (*eventomodels.EventoIngresso, error)
	BuscarItem(usuarioID uint, itemID uint64) (*carrinhomodels.CarrinhoItem, error)
	SalvarItem(item *carrinhomodels.CarrinhoItem) error
	RemoverItem(usuarioID uint, itemID uint64) error
	Limpar(usuarioID uint) error
}

type CarrinhoRepository struct{ db *gorm.DB }

func NewCarrinhoRepository(db *gorm.DB) CarrinhoRepositoryInterface {
	return &CarrinhoRepository{db: db}
}

func (r *CarrinhoRepository) BuscarOuCriar(usuarioID uint) (*carrinhomodels.Carrinho, error) {
	var carrinho carrinhomodels.Carrinho
	err := r.db.Where("usuario_id = ?", usuarioID).First(&carrinho).Error
	if err == nil {
		return &carrinho, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	carrinho = carrinhomodels.Carrinho{UsuarioID: usuarioID}
	if err := r.db.Create(&carrinho).Error; err != nil {
		if err := r.db.Where("usuario_id = ?", usuarioID).First(&carrinho).Error; err != nil {
			return nil, err
		}
	}
	return &carrinho, nil
}

func (r *CarrinhoRepository) Buscar(usuarioID uint) (*carrinhomodels.Carrinho, error) {
	var carrinho carrinhomodels.Carrinho
	err := r.db.Where("usuario_id = ?", usuarioID).
		Preload("Itens.EventoIngresso.Evento").
		Preload("Itens.EventoIngresso.TipoIngresso").
		First(&carrinho).Error
	return &carrinho, err
}

func (r *CarrinhoRepository) BuscarLote(id uint64) (*eventomodels.EventoIngresso, error) {
	var lote eventomodels.EventoIngresso
	err := r.db.Preload("Evento").Preload("TipoIngresso").Preload("Ingressos").
		First(&lote, id).Error
	return &lote, err
}

func (r *CarrinhoRepository) BuscarItem(usuarioID uint, itemID uint64) (*carrinhomodels.CarrinhoItem, error) {
	var item carrinhomodels.CarrinhoItem
	err := r.db.Joins("JOIN tb_carrinho ON tb_carrinho.id = tb_carrinho_item.carrinho_id").
		Where("tb_carrinho.usuario_id = ? AND tb_carrinho_item.id = ?", usuarioID, itemID).
		Preload("EventoIngresso.Evento").Preload("EventoIngresso.TipoIngresso").First(&item).Error
	return &item, err
}

func (r *CarrinhoRepository) SalvarItem(item *carrinhomodels.CarrinhoItem) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var lote eventomodels.EventoIngresso
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Preload("Ingressos").First(&lote, item.EventoIngressoID).Error; err != nil {
			return err
		}
		return tx.Save(item).Error
	})
}

func (r *CarrinhoRepository) RemoverItem(usuarioID uint, itemID uint64) error {
	return r.db.Where("id = ? AND carrinho_id IN (?)", itemID,
		r.db.Model(&carrinhomodels.Carrinho{}).Select("id").Where("usuario_id = ?", usuarioID)).
		Delete(&carrinhomodels.CarrinhoItem{}).Error
}

func (r *CarrinhoRepository) Limpar(usuarioID uint) error {
	return r.db.Where("carrinho_id IN (?)",
		r.db.Model(&carrinhomodels.Carrinho{}).Select("id").Where("usuario_id = ?", usuarioID)).
		Delete(&carrinhomodels.CarrinhoItem{}).Error
}
