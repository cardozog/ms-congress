package repositories

import (
	"strings"
	"time"

	eventomodels "ms-congress/internal/models/evento"

	"gorm.io/gorm"
)

type CheckinRepository struct{ db *gorm.DB }

type CheckinRecord struct {
	IngressoID     uint64     `json:"ingressoId"`
	Codigo         string     `json:"codigo"`
	CheckinEm      *time.Time `json:"checkinEm"`
	Tipo           string     `json:"tipoIngresso"`
	Evento         string     `json:"evento"`
	EventoLogo     string     `json:"eventoLogo"`
	ParticipanteID uint       `json:"participanteId"`
	Nome           string     `json:"nome"`
	CPF            string     `json:"cpf"`
	Email          string     `json:"email"`
	DataNascimento string     `json:"dataNascimento"`
}

func NewCheckinRepository(db *gorm.DB) *CheckinRepository {
	return &CheckinRepository{db: db}
}

func (r *CheckinRepository) Listar(eventoID uint64, organizadorUsuarioID uint, search string) ([]CheckinRecord, error) {
	var records []CheckinRecord
	query := r.db.Table("ingressos AS i").
		Select("i.id AS ingresso_id, i.codigo, i.checkin_em, ti.nome AS tipo, e.nome AS evento, e.evento_logo, u.id AS participante_id, pf.nome, pf.cpf, u.email, pf.data_nascimento").
		Joins("JOIN evento_ingressos AS ei ON ei.id = i.evento_ingresso_id").
		Joins("JOIN eventos AS e ON e.id = ei.evento_id").
		Joins("JOIN tb_pessoa_juridica AS pj ON pj.id = e.organizador_id").
		Joins("JOIN tb_pedido AS p ON p.id = i.pedido_id").
		Joins("JOIN tb_pagamento AS pg ON pg.pedido_id = p.id").
		Joins("JOIN tb_usuario AS u ON u.id = p.usuario_id").
		Joins("JOIN tb_pessoa_fisica AS pf ON pf.usuario_id = u.id").
		Joins("JOIN tipo_ingressos AS ti ON ti.id = ei.tipo_ingresso_id").
		Where("e.id = ? AND pj.usuario_id = ? AND i.status = ? AND p.status = ? AND pg.status = ?", eventoID, organizadorUsuarioID, eventomodels.StatusIngressoCheckin, "PAID", "APPROVED").
		Order("i.checkin_em DESC, i.id DESC")
	search = strings.TrimSpace(search)
	if search != "" {
		like := "%" + search + "%"
		query = query.Where("pf.nome ILIKE ? OR pf.cpf ILIKE ? OR u.email ILIKE ?", like, like, like)
	}
	return records, query.Scan(&records).Error
}

func (r *CheckinRepository) EventoPertenceAoOrganizador(eventoID uint64, organizadorUsuarioID uint) (bool, error) {
	var count int64
	err := r.db.Table("eventos AS e").
		Joins("JOIN tb_pessoa_juridica AS pj ON pj.id = e.organizador_id").
		Where("e.id = ? AND pj.usuario_id = ?", eventoID, organizadorUsuarioID).
		Count(&count).Error
	return count > 0, err
}
