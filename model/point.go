package model

import (
	"gorm.io/gorm"
)

type Point struct {
	*gorm.Model
	Name   string `json:"name" form:"name"`
	UserId int    `json:"userId" form:"userId"`
	Amount int    `json:"amount" form:"amount"`
}
