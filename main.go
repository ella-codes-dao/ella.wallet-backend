package main

import (
	"fmt"

	"ella.wallet-backend/internal/config"
	"ella.wallet-backend/internal/controllers"
	"ella.wallet-backend/internal/database/migrate"
	"ella.wallet-backend/internal/services"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load Configurations using Viper
	config.LoadAppConfig()

	// Initialize DB
	services.StartDB()
	migrate.Migrate(services.DB)

	// Initialize Services
	services.Start()

	// Initialize Router
	router := initRouter()
	router.Run(fmt.Sprintf(":%v", config.AppConfig.ServerPort))
}

func initRouter() *gin.Engine {
	router := gin.Default()

	network := router.Group("/:network")
	{
		// Submit public key to receive accounts
		network.GET("/address", controllers.GetAddress)

		// Submit a public key to create a new account
		network.POST("/address", controllers.CreateAddress)
	}

	return router
}
