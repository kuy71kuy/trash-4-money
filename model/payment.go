package model

import "gorm.io/gorm"

type Payment struct {
	*gorm.Model
	UserId      int    `json:"user_id" form:"user_id"`
	PointId     int    `json:"point_id" form:"point_id"`
	Amount      int    `json:"amount" form:"amount"`
	Type        string `json:"type" form:"type"`
	Number      string `json:"number" form:"number"`
	Status      string `json:"status" form:"status"`
	ReferenceNo string `json:"reference_no" form:"reference_no"`
}
