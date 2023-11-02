package model

import "gorm.io/gorm"

type Payment struct {
	*gorm.Model
	UserId  int    `json:"userId" form:"userId"`
	PointId int    `json:"pointId" form:"pointId"`
	Amount  int    `json:"amount" form:"amount"`
	Type    string `json:"type" form:"type"`
	Number  string `json:"number" form:"number"`
	Status  string `json:"status" form:"status"`
}
