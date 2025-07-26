package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port        int
	DatabaseURL string
	JWTSecret   string
	Environment string
}

func Load() *Config {
	return &Config{
		Port:        getEnvAsInt("PORT", 8080),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/dbname"),
		JWTSecret:   getEnv("JWT_SECRET", "defaultsecret"),
		Environment: getEnv("ENVIRONMENT", "development"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
