package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WorkspaceType string

var (
	Private WorkspaceType = "private"
	Public  WorkspaceType = "public"
)

func (wt WorkspaceType) String() string {
	if wt == Private ||
		wt == Public {
		return string(wt)
	}
	return ""
}

type Workspace struct {
	ID           string                 `gorm:"column:workspace_id;type:uuid;default:uuid_generate_v4();primaryKey" json:"workspaceId"`
	Type         WorkspaceType          `gorm:"column:workspace_type" json:"workspaceType"`
	Name         string                 `gorm:"column:workspace_name" json:"workspaceName"`
	DisplayName  string                 `gorm:"column:workspace_display_name" json:"workspaceDisplayName"`
	Environments []WorkspaceEnvironment `gorm:"->;foreignKey:WorkspaceID" json:"environments"`
	Collections  Collections            `gorm:"->;foreignKey:WorkspaceID" json:"collections"`
	CreatedAt    time.Time              `gorm:"column:created_at" json:"-"`
	UpdatedAt    time.Time              `gorm:"column:updated_at" json:"-"`
}

func (m *Workspace) BeforeCreate(tx *gorm.DB) (err error) {
	current := time.Now()
	m.ID = uuid.NewString()
	m.Type = Private
	m.CreatedAt = current
	m.UpdatedAt = current
	return
}

func (m *Workspace) BeforeUpdate(tx *gorm.DB) (err error) {
	m.UpdatedAt = time.Now()
	return
}
