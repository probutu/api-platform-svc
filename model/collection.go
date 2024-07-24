package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Collection struct {
	ID          string         `gorm:"column:collection_id;type:uuid;default:uuid_generate_v4();primaryKey" json:"collectionId"`
	WorkspaceID string         `gorm:"column:workspace_id;type:uuid" form:"workspaceId" json:"workspaceId"`
	Name        string         `gorm:"column:collection_name" json:"collectionName"`
	Folders     []Folder       `gorm:"->;foreignKey:CollectionID" json:"folders"`
	Requests    []*Request     `gorm:"->;foreignKey:CollectionID" json:"requests"`
	CreatedAt   time.Time      `gorm:"column:created_at" json:"-"`
	UpdatedAt   time.Time      `gorm:"column:updated_at" json:"-"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at" json:"-"`
}

type Collections []Collection

func (m *Collection) BeforeCreate(tx *gorm.DB) (err error) {
	current := time.Now()
	m.ID = uuid.NewString()
	m.CreatedAt = current
	m.UpdatedAt = current
	return
}

func (m *Collection) BeforeUpdate(tx *gorm.DB) (err error) {
	m.UpdatedAt = time.Now()
	return
}
