package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

type Config struct {
	Port        int
	DatabaseURL string
	JWTSecret   string
}

func LoadConfig() *Config {
	cfg := &Config{
		Port:        getEnvAsInt("PORT", 8080),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		JWTSecret:   getEnv("JWT_SECRET", "change-me-in-prod"),
	}

	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	return cfg
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	fmt.Println("Getting DB env")
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		fmt.Println(value)
		return value
	}
	return defaultVal
}
