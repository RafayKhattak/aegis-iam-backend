package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

const defaultPort = "8080"

// AppConfig contains runtime configuration values for the API.
type AppConfig struct {
	Port string
}

// LoadConfig loads application configuration from the environment.
//
// It first attempts to load values from a local .env file for development.
// If the file is not present (common in containers/production), it falls back
// to existing system environment variables.
func LoadConfig() AppConfig {
	if err := loadDotEnv(); err != nil && !os.IsNotExist(err) {
		log.Printf("config: %v; continuing with system environment variables", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	return AppConfig{
		Port: port,
	}
}

func loadDotEnv() error {
	return godotenv.Load()
}
