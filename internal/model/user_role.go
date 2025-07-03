package model

type UserRole struct {
	BaseModel
	Code UserRoleEnum `gorm:"not null;unique" json:"code"`
	Name string       `gorm:"not null" json:"name"`
}

type UserRoleEnum string

const (
	UserRoleCodeSuperAdmin UserRoleEnum = "role:super_admin"
	UserRoleCodeAdmin      UserRoleEnum = "role:admin"
	UserRoleCodeUser       UserRoleEnum = "role:user"
)

var UserRoleSet = []UserRole{
	{
		Code: UserRoleCodeSuperAdmin,
		Name: "超级管理员",
	},
	{
		Code: UserRoleCodeAdmin,
		Name: "管理员",
	},
	{
		Code: UserRoleCodeUser,
		Name: "用户",
	},
}
