package model

import (
	"time"

	"gorm.io/gorm"
)

const (
	AdminRole  string = "admin"
	UserRole   string = "user"
	DummyToken string = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwibmFtZSI6IiIsInJvbGUiOiIiLCJleHAiOjEwMDAwMDAwMDB9.VMb4PNNJ_CUq7oom1P1kx5_6VSLOZIP1ZVg9T4rJIMA"
)

type User struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string         `json:"name" form:"name"`
	Email     string         `json:"email" form:"email"`
	Password  string         `json:"password" form:"password"`
	Role      string         `json:"role" form:"role"`
}
