package model

type User struct {
	BaseModel
	UniqueID  int64       `gorm:"index" json:"uniqueId"`
	Email     string      `gorm:"index" json:"email"`
	Mobile    string      `gorm:"index" json:"mobile"`
	Password  string      `json:"-"`
	Salt      string      `json:"-"`
	Nickname  string      `gorm:"index" json:"nickname"`
	AvatarURL string      `json:"avatarUrl"`
	Roles     []*UserRole `gorm:"many2many:user_role_ref;" json:"roles"`
}

func (u *User) TableName() string {
	return "users"
}
