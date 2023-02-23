package database

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connect(connectionString string) (*gorm.DB, error) {
	instance, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	log.Println("Connected to Database!")

	return instance, nil
}
