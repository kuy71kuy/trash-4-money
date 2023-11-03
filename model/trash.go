package model

import "gorm.io/gorm"

type Trash struct {
	*gorm.Model
	UserId  int    `json:"user_id" form:"user_id"`
	Type    string `json:"type" form:"type"`
	Weight  int    `json:"weight" form:"weight"`
	Address string `json:"address" form:"address"`
	Image   string `json:"image" form:"image"`
	Note    string `json:"note" form:"note"`
	Status  string `json:"status" form:"status"`
}
