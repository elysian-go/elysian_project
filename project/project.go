package project

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Base struct {
	ID			string		`json:"_id,omitempty", bson:"_id,omitempty"`
	DeletedAt	*time.Time	`sql:"index"`
	CreatedAt	time.Time
	UpdatedAt	time.Time
}

type Project struct {
	Base
	Title				string		`json:"title" bson:"title,omitempty"`
	Description			string		`json:"description" bson:"description,omitempty"`
	Owner				string		`json:"owner" bson:"owner"`
	Archived			bool		`json:"archived" bson:"archived"`

}

func (m *Base) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("UpdatedAt", time.Now().Unix())
	return nil
}

func (m *Base) BeforeCreate(scope *gorm.Scope) error {
	if m.UpdatedAt.IsZero() {
		scope.SetColumn("UpdatedAt", time.Now().Unix())
	}

	scope.SetColumn("CreatedAt", time.Now().Unix())
	return nil
}