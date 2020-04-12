package project

type OwnerProject struct {
	ID          uint64 `gorm:"primary_key"`
	UserId		string `gorm:"type:uuid;not_null;index"`
	ProjectId	string	`gorm:"type:char(24)"`
}
