package project

import "time"

type Model struct {
	ID					string		`json:"id,omitempty"`
	Title				string		`json:"title" binding:"required,max=15"`
	Description			string		`json:"description" binding:"required,max=512"`
	Owner				string		`json:"owner,omitempty"` //TODO figure what to return
	Archived			bool		`json:"archived,omitempty"`
	CreatedAt			time.Time	`json:"created_at,omitempty"`
	UpdatedAt			time.Time	`json:"updated_at,omitempty"`
}

type AddContributorModel struct {
	ID []string `json:"ids"`
}