package dto

type UserLoginByEmailResDTO struct {
	Token     string `json:"token"`
	RefreshAt int64  `json:"refreshAt"`
	ExpireAt  int64  `json:"expireAt"`
}
