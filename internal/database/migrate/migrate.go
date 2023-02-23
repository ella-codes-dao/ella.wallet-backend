package migrate

import (
	"log"

	"ella.wallet-backend/internal/models"
	"gorm.io/gorm"
)

func Migrate(instance *gorm.DB) {
	instance.AutoMigrate(
		models.Wallet{},
	)

	log.Println("Database Migration Completed!")
}
