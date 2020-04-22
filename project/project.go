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

type OwnerProject struct {
	UserId		string `gorm:"primary_key;type:uuid;not_null;index:idx_owner_project"`
	ProjectId	string	`gorm:"primary_key;type:char(24);unique_index:idx_owner_project"`
}

type CollaboratorProject struct {
	UserId		string `gorm:"primary_key;type:uuid;not_null;index:idx_collab_project"`
	ProjectId	string	`gorm:"primary_key;type:char(24);index:idx_collab_project"`
}
