package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/lib/pq"
)

type Entities struct {
	Indices pq.Int32Array `gorm:"column:indices" json:"indices"`
	Text    string        `gorm:"column:text" json:"text"`
}

type Query struct {
	Entities
}

func (a Query) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *Query) Scan(value interface{}) error {
	b, ok := value.(string)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal([]byte(b), &a)
}

type Path struct {
	Entities
}

func (a Path) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *Path) Scan(value interface{}) error {
	b, ok := value.(string)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal([]byte(b), &a)
}

type Param struct {
	Query []Query `gorm:"column:query;type:VARCHAR(255)" json:"queries"`
	Path  []Path  `gorm:"column:path;type:VARCHAR(255)" json:"paths"`
}

func (a Param) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *Param) Scan(value interface{}) error {
	b, ok := value.(string)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal([]byte(b), &a)
}

type Body map[string]interface{}

func (a Body) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *Body) Scan(value interface{}) error {
	b, ok := value.(string)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal([]byte(b), &a)
}
