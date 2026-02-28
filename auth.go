package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL     string
	JWTSecret       string
	Port            string
	RateLimitPerMin int
}

var C Config

func Load() {
	if err := godotenv.Load(); err != nil {
		log.Println("[Config] .env file not found, reading from environment")
	}

	C = Config{
		DatabaseURL:     mustGet("DATABASE_URL"),
		JWTSecret:       mustGet("JWT_SECRET"),
		Port:            getOrDefault("PORT", "8080"),
		RateLimitPerMin: getIntOrDefault("RATE_LIMIT_PER_MIN", 60),
	}
}

func mustGet(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("[Config] Required env var %s is not set", key)
	}
	return val
}

func getOrDefault(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getIntOrDefault(key string, fallback int) int {
	if val := os.Getenv(key); val != "" {
		if n, err := strconv.Atoi(val); err == nil {
			return n
		}
	}
	return fallback
}
