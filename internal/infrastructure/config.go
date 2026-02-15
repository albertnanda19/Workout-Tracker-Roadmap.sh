package infrastructure

import (
	"errors"
	"os"
)

type Config struct {
	Port        string
	DatabaseURL string
}

func LoadConfig() (Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return Config{}, errors.New("DATABASE_URL is required")
	}

	return Config{
		Port:        port,
		DatabaseURL: databaseURL,
	}, nil
}
