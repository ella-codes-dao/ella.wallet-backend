package services

import (
	"log"

	"ella.wallet-backend/internal/config"
	"ella.wallet-backend/internal/database"
	"ella.wallet-backend/internal/flow"
	"ella.wallet-backend/internal/queue"
	"gorm.io/gorm"
)

var DB *gorm.DB
var Flow *flow.FlowClient
var Queues *queue.JobQueues

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
	client, err := flow.NewClient()
	if err != nil {
		log.Fatalf("Unable to start Flow Client: %s", err)
	}

	Flow = client
}

func startQueues() {
	Queues = queue.StartQueues()

	go Queues.KuCoin.DoWork()
}
