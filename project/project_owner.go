package project

type OwnerProject struct {
	ID          uint64 `gorm:"primary_key"`
	UserId		string `gorm:"type:uuid;not_null;index:idx_user_project"`
	ProjectId	string	`gorm:"type:char(24);unique_index:idx_user_project"`
}
