package model

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RequestMethod string

var (
	GET     RequestMethod = http.MethodGet
	POST    RequestMethod = http.MethodPost
	PUT     RequestMethod = http.MethodPut
	PATCH   RequestMethod = http.MethodPatch
	DELETE  RequestMethod = http.MethodDelete
	HEAD    RequestMethod = http.MethodHead
	OPTIONS RequestMethod = http.MethodOptions
)

type Request struct {
	ID           string     `gorm:"column:request_id;type:uuid;default:uuid_generate_v4();primaryKey;" json:"requestId"`
	WorkspaceID  string     `gorm:"column:workspace_id;type:uuid" form:"workspaceId" json:"workspaceId"`
	CollectionID *string    `gorm:"column:collection_id;type:uuid;" form:"collectionId" json:"collectionId"`
	FolderID     *string    `gorm:"column:folder_id;type:uuid;" form:"folderId" json:"folderId"`
	Name         string     `gorm:"column:name" json:"name"`
	Method       string     `gorm:"column:method" json:"method"`
	URI          string     `gorm:"column:uri" json:"uri"`
	Params       Param      `gorm:"column:params;type:VARCHAR(255)" json:"params"`
	Headers      []Header   `gorm:"foreignKey:RequestID" json:"headers"`
	Body         Body       `gorm:"column:body;type:VARCHAR(255)" json:"body"`
	Responses    []Response `gorm:"foreignKey:RequestID" json:"response"`
	CreatedAt    time.Time  `gorm:"column:created_at" json:"-"`
	UpdatedAt    time.Time  `gorm:"column:updated_at" json:"-"`
}

type Requests []Request

func (m *Request) BeforeCreate(tx *gorm.DB) (err error) {
	current := time.Now()
	m.ID = uuid.NewString()
	m.CreatedAt = current
	m.UpdatedAt = current
	return
}

func (m *Request) BeforeUpdate(tx *gorm.DB) (err error) {
	m.UpdatedAt = time.Now()
	return
}
