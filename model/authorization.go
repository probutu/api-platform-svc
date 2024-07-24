package model

import "time"

type AuthorizationType string

var (
	Basic  AuthorizationType = "basic"
	Bearer AuthorizationType = "bearer"
	Jwt    AuthorizationType = "jwt"
	OAuth1 AuthorizationType = "oauth1"
	OAuth2 AuthorizationType = "oauth2"
)

type Authorization struct {
	ID        string            `gorm:"column:authorization_id;type:uuid;default:uuid_generate_v4();primaryKey"`
	Type      AuthorizationType `gorm:"column:authorization_type" json:"authorization_type"`
	CreatedAt time.Time         `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time         `gorm:"column:updated_at" json:"updatedAt"`
}
