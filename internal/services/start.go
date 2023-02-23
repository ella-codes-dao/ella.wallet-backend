package services

import (
	"log"

	"ella.wallet-backend/internal/config"
	"ella.wallet-backend/internal/database"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Start() {
	startFlow()
	startQueues()
}

func StartDB() {
	instance, err := database.Connect(config.AppConfig.GormConnection)
	if err != nil {
		log.Fatalf("Unable to connect to database: %s", err)
	}

	DB = instance
}

func startFlow() {

}

func startQueues() {

}
