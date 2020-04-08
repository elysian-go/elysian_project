package project

import "time"

type Model struct {
	ID					string		`json:"id,string"`
	Title				string		`json:"title"`
	Description			string		`json:"description"`
	Owner				string		`json:"owner"`
	Archived			bool		`json:"archived"`
	CreatedAt			time.Time	`json:"created_at,omitempty"`
	UpdatedAt			time.Time	`json:"updated_at,omitempty"`
}