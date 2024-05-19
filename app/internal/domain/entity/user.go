package entity

type User struct {
	ID       int64
	Username string
	Email    string
	Password string
}

type UserInfo struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UserView struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

type OAuthUserDTO struct {
	Email string
}

type LoginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterUserDTO struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type SetUsernameDTO struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
}
