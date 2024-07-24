package model

import (
	"time"

	"gorm.io/gorm"
)

type Header struct {
	ID           string    `gorm:"column:header_id;type:uuid;default:uuid_generate_v4();primaryKey" json:"-"`
	RequestID    string    `gorm:"column:request_id;type:uuid" json:"-"`
	Key          string    `gorm:"column:key" json:"key"`
	Value        string    `gorm:"column:value" json:"value"`
	Description  string    `gorm:"column:description" json:"description"`
	IsCanDeleted bool      `gorm:"column:is_can_deleted" json:"isCanDeleted"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (m *Header) BeforeCreate(tx *gorm.DB) (err error) {
	current := time.Now()
	m.CreatedAt = current
	m.UpdatedAt = current
	return
}

func (m *Header) BeforeUpdate(tx *gorm.DB) (err error) {
	m.UpdatedAt = time.Now()
	return
}
