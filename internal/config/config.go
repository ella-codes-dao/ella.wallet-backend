package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	GormConnection  string `mapstructure:"GORM_CONNECTION"`
	ServerPort      string `mapstructure:"SERVER_PORT"`
	FlowNetwork     string `mapstructure:"FLOW_NETWORK"`
	FlowBackendAddr string `mapstructure:"FLOW_BACKEND_ADDR"`
	FlowBackendKey  string `mapstructure:"FLOW_BACKEND_KEY"`
	FCLAppID        string `mapstructure:"FCL_APP_ID"`
	GcpProjectID    string `mapstructure:"GCP_PROJECT_ID"`
}

var AppConfig *Config

func LoadAppConfig() {
	log.Println("Loading Server Configurations...")
	// Read file path
	viper.AddConfigPath(".")
	// set config file and path
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	// watching changes in app.env
	viper.AutomaticEnv()
	// reading the config file
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = viper.Unmarshal(&AppConfig)
	if err != nil {
		log.Fatal(err)
	}
}
