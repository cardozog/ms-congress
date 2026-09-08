package repositories

import (
	models "ms-congress/internal/models/evento"

	"gorm.io/gorm"
)

type EventoRepository struct {
	db *gorm.DB
}

type EventoRepositoryInterface interface {
	// Define os métodos que o repositório de eventos deve implementar
}

func NewEventoRepository(db *gorm.DB) EventoRepositoryInterface {
	return &EventoRepository{
		db: db,
	}
}

func (r *EventoRepository) Create(evento *models.Evento) error {
	if err := r.db.Create(evento).Error; err != nil {
		return err
	}
	return nil
}

func (r *EventoRepository) GetByID(id uint64) (*models.Evento, error) {
	var evento models.Evento
	if err := r.db.First(&evento, id).Error; err != nil {
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
