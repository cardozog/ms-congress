package repositories

import (
	"time"

	models "ms-congress/internal/models/evento"

	"gorm.io/gorm"
)

type EventoRepository struct {
	db *gorm.DB
}

type EventoRepositoryInterface interface {
	Create(evento *models.Evento) error
	GetByID(id uint64) (*models.Evento, error)
	BuscarEventosPorOrganizador(organizadorID uint64) ([]models.Evento, error)
	Update(evento *models.Evento) error
	Publicar(id uint64) error
	Encerrar(id uint64, dataFim time.Time) error
	Delete(evento *models.Evento) error
	ListarTiposIngresso() ([]models.TipoIngresso, error)
}

func NewEventoRepository(db *gorm.DB) EventoRepositoryInterface {
	return &EventoRepository{
		db: db,
	}
}

func (r *EventoRepository) Create(evento *models.Evento) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		ingressos := evento.Ingressos
		evento.Ingressos = nil
		if err := tx.Create(evento).Error; err != nil {
			return err
		}
		for index := range ingressos {
			ingresso := &ingressos[index]
			ingresso.EventoID = evento.ID
			if err := tx.Create(ingresso).Error; err != nil {
				return err
			}
		}
		evento.Ingressos = ingressos
		return nil
	})
}

func (r *EventoRepository) ListarTiposIngresso() ([]models.TipoIngresso, error) {
	var tipos []models.TipoIngresso
	if err := r.db.Order("id ASC").Find(&tipos).Error; err != nil {
		return nil, err
	}
	return tipos, nil
}

func (r *EventoRepository) GetByID(id uint64) (*models.Evento, error) {
	var evento models.Evento
	if err := r.db.
		Preload("Organizador").
		Preload("Ingressos").
		Preload("Ingressos.TipoIngresso").
		Preload("Ingressos.Ingressos").
		First(&evento, id).Error; err != nil {
		return nil, err
	}
	return &evento, nil
}

func (r *EventoRepository) BuscarEventosPorOrganizador(organizadorID uint64) ([]models.Evento, error) {
	var eventos []models.Evento
	if err := r.db.Where("organizador_id = ?", organizadorID).Find(&eventos).Error; err != nil {
		return nil, err
	}
	return eventos, nil
}

func (r *EventoRepository) Update(evento *models.Evento) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Evento{}).Where("id = ?", evento.ID).Updates(map[string]interface{}{
			"nome":           evento.Nome,
			"descricao":      evento.Descricao,
			"data_inicio":    evento.DataInicio,
			"data_fim":       evento.DataFim,
			"publicado":      evento.Publicado,
			"organizador_id": evento.OrganizadorID,
			"logradouro":     evento.Logradouro,
			"cidade":         evento.Cidade,
			"estado":         evento.Estado,
			"cep":            evento.CEP,
			"evento_logo":    evento.EventoLogo,
		}).Error; err != nil {
			return err
		}

		for _, ingresso := range evento.Ingressos {
			updates := map[string]interface{}{
				"preco":             ingresso.Preco,
				"quantidade":        ingresso.Quantidade,
				"data_inicio_venda": ingresso.DataInicioVenda,
				"data_fim_venda":    ingresso.DataFimVenda,
			}
			result := tx.Model(&models.EventoIngresso{}).
				Where("evento_id = ? AND tipo_ingresso_id = ?", evento.ID, ingresso.TipoIngressoID).
				Updates(updates)
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				if err := tx.Create(&models.EventoIngresso{
					EventoID: evento.ID, TipoIngressoID: ingresso.TipoIngressoID,
					Preco: ingresso.Preco, Quantidade: ingresso.Quantidade,
					DataInicioVenda: ingresso.DataInicioVenda, DataFimVenda: ingresso.DataFimVenda,
				}).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (r *EventoRepository) Publicar(id uint64) error {
	return r.db.Model(&models.Evento{}).Where("id = ?", id).Update("publicado", true).Error
}

func (r *EventoRepository) Encerrar(id uint64, dataFim time.Time) error {
	return r.db.Model(&models.Evento{}).Where("id = ?", id).Update("data_fim", dataFim).Error
}

func (r *EventoRepository) Delete(evento *models.Evento) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("evento_ingresso_id IN (?)",
			tx.Model(&models.EventoIngresso{}).
				Select("id").
				Where("evento_id = ?", evento.ID),
		).Delete(&models.Ingresso{}).Error; err != nil {
			return err
		}
		if err := tx.Where("evento_id = ?", evento.ID).
			Delete(&models.EventoIngresso{}).Error; err != nil {
			return err
		}
		if err := tx.Where("id = ?", evento.ID).
			Delete(&models.Evento{}).Error; err != nil {
			return err
		}
		return nil
	})
}
