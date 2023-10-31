package config

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"project/model"
)

var (
	DB *gorm.DB
)

func initDB() {
	config := model.Config{
		DbUsername: os.Getenv("DB_USERNAME"),
		DbPassword: os.Getenv("DB_PASSWORD"),
		DbHost:     os.Getenv("DB_HOST"),
		DbPort:     os.Getenv("DB_PORT"),
		DbName:     os.Getenv("DB_NAME"),
	}
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		config.DbUsername,
		config.DbPassword,
		config.DbHost,
		config.DbPort,
		config.DbName,
	)
	var err error
	DB, err = gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	initialMigration()
}

func initialMigration() {
	DB.AutoMigrate(&model.User{})
}
