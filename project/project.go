package project

import (
	"time"
)

type Base struct {
	ID			string		`json:"-" bson:"_id,omitempty"`
	DeletedAt	*time.Time	`sql:"index"`
	CreatedAt	time.Time
	UpdatedAt	time.Time
}

type Project struct {
	Base							`bson:",inline"`
	Title				string		`json:"title" bson:"title,omitempty"`
	Description			string		`json:"description" bson:"description,omitempty"`
	Owner				string		`json:"owner" bson:"owner"`
	Archived			bool		`json:"archived" bson:"archived"`
}