package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	AppName string
	AppEnv  string
}

var GlobalConfig AppConfig

func (cfg *AppConfig) LoadVariables(envPath ...string) error {
	err := godotenv.Load(envPath...)

	if err != nil {
		log.Println(".env file not found. Loading from system environment", err)
	}

	cfg.AppName = os.Getenv("APP_NAME")
	cfg.AppEnv = os.Getenv("APP_ENV")

	return nil
}
