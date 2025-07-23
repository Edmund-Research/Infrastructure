package config

import (
	"os"
)

type Config struct {
	Port         string
	MongoURI     string
	DatabaseName string
	SOARWebhook  string
}

func Load() (*Config, error) {
	cfg := &Config{
		Port:         getEnv("PORT", "8080"),
		MongoURI:     getEnv("MONGO_URI", "mongodb://localhost:27017"),
		DatabaseName: getEnv("DATABASE_NAME", "threat_intelligence"),
		SOARWebhook:  getEnv("SOAR_WEBHOOK", "https://soar.example.com/webhook"),
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
