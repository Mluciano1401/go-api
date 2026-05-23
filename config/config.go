package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort    string
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	JWTSecret     string
	JWTExpiresHrs int
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Aviso: no se encontró archivo .env, usando variables del entorno")
	}

	expHrs, _ := strconv.Atoi(getEnv("JWT_EXPIRES_HOURS", "24"))

	return &Config{
		ServerPort:    getEnv("SERVER_PORT", "8080"),
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBPort:        getEnv("DB_PORT", "5433"),
		DBUser:        getEnv("DB_USER", "postgres"),
		DBPassword:    getEnv("DB_PASSWORD", "postgres"),
		DBName:        getEnv("DB_NAME", "tarea1_db"),
		JWTSecret:     getEnv("JWT_SECRET", "1234567890"),
		JWTExpiresHrs: expHrs,
	}
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}
