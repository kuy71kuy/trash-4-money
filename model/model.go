package model

import "gorm.io/gorm"

type Config struct {
	DbUsername string
	DbPassword string
	DbPort     string
	DbHost     string
	DbName     string
}

type User struct {
	*gorm.Model
	Name     string `json:"name" form:"name"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
	Role     string `json:"role" form:"role"`
}
