package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Folder struct {
	ID           string     `gorm:"column:folder_id;type:uuid;default:uuid_generate_v4();primaryKey" json:"-"`
	CollectionID string     `gorm:"column:collection_id;type:uuid" form:"collectionId" json:"collectionId"`
	Name         string     `gorm:"column:name" form:"name" json:"name"`
	Requests     []*Request `gorm:"->;foreignKey:FolderID" json:"requests"`
	CreatedAt    time.Time  `gorm:"column:created_at" json:"-"`
	UpdatedAt    time.Time  `gorm:"column:updated_at" json:"-"`
}

func (m *Folder) BeforeCreate(tx *gorm.DB) (err error) {
	current := time.Now()
	m.ID = uuid.NewString()
	m.CreatedAt = current
	m.UpdatedAt = current
	return
}

func (m *Folder) BeforeUpdate(tx *gorm.DB) (err error) {
	m.UpdatedAt = time.Now()
	return
}

type Folders []Folder
