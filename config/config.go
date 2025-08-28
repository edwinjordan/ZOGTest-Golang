package config

import (
	"context"

	"github.com/edwinjordan/ZOGTest-Golang.git/internal/logging"
	"github.com/joho/godotenv"
)

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		logging.LogInfo(context.Background(), "No .env file found, using environment variables")

	}
}
