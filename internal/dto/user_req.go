package dto

type UserLoginByEmailReqDTO struct {
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required"`
}
