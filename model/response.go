package model

import (
	"time"

	"gorm.io/gorm"
)

type Response struct {
	ID        string    `gorm:"column:response_id;type:uuid;default:uuid_generate_v4();primaryKey" json:"responseId"`
	RequestID string    `gorm:"column:request_id;type:uuid" json:"requestId"`
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (m *Response) BeforeCreate(tx *gorm.DB) (err error) {
	current := time.Now()
	m.CreatedAt = current
	m.UpdatedAt = current
	return
}

func (m *Response) BeforeUpdate(tx *gorm.DB) (err error) {
	m.UpdatedAt = time.Now()
	return
}
