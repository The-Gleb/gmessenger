package dto

type RegisterUserDTO struct {
	UserName string `json:"username"`
	Login    string `json:"login"`
	Password string `json:"password"`
}
