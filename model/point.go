package model

import (
	"gorm.io/gorm"
)

type Point struct {
	*gorm.Model
	Name   string `json:"name" form:"name"`
	UserId int    `json:"user_id" form:"user_id"`
	Amount int    `json:"amount" form:"amount"`
}
