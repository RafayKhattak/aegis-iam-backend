package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

const defaultPort = "8080"
const defaultTokenDuration = 24 * time.Hour

// AppConfig contains runtime configuration values for the API.
type AppConfig struct {
	Port          string
	DBSource      string
	JWTSecret     string
	TokenDuration time.Duration
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

	tokenDuration := defaultTokenDuration
	if durationStr := os.Getenv("TOKEN_DURATION"); durationStr != "" {
		parsedDuration, err := time.ParseDuration(durationStr)
		if err != nil {
			log.Printf("config: invalid TOKEN_DURATION %q; using default %s", durationStr, defaultTokenDuration)
		} else {
			tokenDuration = parsedDuration
		}
	}

	return AppConfig{
		Port:          port,
		DBSource:      os.Getenv("DB_SOURCE"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
		TokenDuration: tokenDuration,
	}
}

func loadDotEnv() error {
	return godotenv.Load()
}
