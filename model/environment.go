package model

import "time"

type WorkspaceEnvironment struct {
	ID           string        `gorm:"column:environment_id;type:uuid;default:uuid_generate_v4();primaryKey" json:"-"`
	WorkspaceID  string        `gorm:"column:workspace_id;type:uuid" json:"-"`
	Name         string        `gorm:"column:environment_name" json:"environment_name"`
	Environments []Environment `gorm:"foreignKey:WorkspaceEnvironmentID" json:"environments"`
}

type Environment struct {
	ID                     string    `gorm:"column:environment_id;type:uuid;default:uuid_generate_v4();primaryKey" json:"-"`
	WorkspaceEnvironmentID string    `gorm:"column:workspace_environment_id;type:uuid" json:"-"`
	Variable               string    `gorm:"column:variable" json:"variable"`
	InitialValue           string    `gorm:"column:initial_value" json:"initialValue"`
	CurrentValue           string    `gorm:"column:current_value" json:"currentValue"`
	CreatedAt              time.Time `gorm:"column:created_at" json:"-"`
	UpdatedAt              time.Time `gorm:"column:updated_at" json:"-"`
}
