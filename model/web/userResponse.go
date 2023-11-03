package web

type UserResponse struct {
	Id       int    `json:"id"`
	Name     string `json:"name" form:"name"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

type GetUserResponse struct {
	Id    int    `json:"id"`
	Name  string `json:"name" form:"name"`
	Email string `json:"email" form:"email"`
}

type UserLoginResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Token string `json:"token"`
}
