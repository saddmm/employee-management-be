package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv             string
	AppPort            string
	DBHost             string
	DBPort             string
	DBUser             string
	DBPassword         string
	DBName             string
	JWTSecret          string
	JWTExpireHours     int
	RateLimitMax       int
	RateLimitWindowSec int
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Note: .env file not found or could not be loaded, reading from environment")
	}

	jwtExpireHours, err := strconv.Atoi(getEnv("JWT_EXPIRE_HOURS", "24"))
	if err != nil {
		jwtExpireHours = 24
	}

	rateLimitMax, err := strconv.Atoi(getEnv("RATE_LIMIT_MAX", "100"))
	if err != nil {
		rateLimitMax = 100
	}

	rateLimitWindowSec, err := strconv.Atoi(getEnv("RATE_LIMIT_WINDOW_SEC", "60"))
	if err != nil {
		rateLimitWindowSec = 60
	}

	return &Config{
		AppEnv:             getEnv("APP_ENV", "development"),
		AppPort:            getEnv("APP_PORT", "8080"),
		DBHost:             getEnv("DB_HOST", "localhost"),
		DBPort:             getEnv("DB_PORT", "3306"),
		DBUser:             getEnv("DB_USER", "ems_user"),
		DBPassword:         getEnv("DB_PASSWORD", "ems_password"),
		DBName:             getEnv("DB_NAME", "employee_management"),
		JWTSecret:          getEnv("JWT_SECRET", "default_secret_key_please_change"),
		JWTExpireHours:     jwtExpireHours,
		RateLimitMax:       rateLimitMax,
		RateLimitWindowSec: rateLimitWindowSec,
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return fallback
}
