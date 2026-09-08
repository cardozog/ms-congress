package repositories

import (
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
	if err := r.db.Save(evento).Error; err != nil {
		return err
	}
	return nil
}

func (r *EventoRepository) Delete(evento *models.Evento) error {
	if err := r.db.Delete(evento).Error; err != nil {
		return err
	}
	return nil
}
