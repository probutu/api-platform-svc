package model

import (
	"time"

	"gorm.io/gorm"
)

type WorkspaceUserRole string

var (
	Admin WorkspaceUserRole = "admin"
	Write WorkspaceUserRole = "write"
	Read  WorkspaceUserRole = "read"
)

type WorkspaceUser struct {
	ID          string              `gorm:"column:user_id;type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID      string              `gorm:"column:user_id;type:uuid" json:"userId"`
	WorkspaceID string              `gorm:"column:workspace_id;type:uuid" json:"workspaceId"`
	Roles       []WorkspaceUserRole `gorm:"column:roles;type:varchar(255)" json:"roles"`
	CreatedAt   time.Time           `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   time.Time           `gorm:"column:updated_at" json:"updatedAt"`
}

func (m *WorkspaceUser) BeforeCreate(tx *gorm.DB) (err error) {
	current := time.Now()
	m.CreatedAt = current
	m.UpdatedAt = current
	return
}

func (m *WorkspaceUser) BeforeUpdate(tx *gorm.DB) (err error) {
	m.UpdatedAt = time.Now()
	return
}
