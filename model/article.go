package model

import "gorm.io/gorm"

type Article struct {
	*gorm.Model
	Title     string `json:"title" form:"title"`
	Text      string `json:"text" form:"text"`
	Link      string `json:"link" form:"link"`
	Thumbnail string `json:"thumbnail" form:"thumbnail"`
}
