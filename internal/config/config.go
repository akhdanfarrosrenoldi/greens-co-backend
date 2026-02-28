package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                    string
	DatabaseURL             string
	JWTSecret               string
	JWTExpiresIn            string
	MidtransServerKey       string
	MidtransClientKey       string
	MidtransProduction      bool
	CORSOrigin              string
	AdminDefaultPassword    string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	jwtSecret := getEnv("JWT_SECRET", "")
	if len(jwtSecret) < 32 {
		log.Fatal("FATAL: JWT_SECRET must be at least 32 characters long. Set a strong secret in .env")
	}

	return &Config{
		Port:                 getEnv("PORT", "8080"),
		DatabaseURL:          getEnv("DATABASE_URL", ""),
		JWTSecret:            jwtSecret,
		JWTExpiresIn:         getEnv("JWT_EXPIRES_IN", "7d"),
		MidtransServerKey:    getEnv("MIDTRANS_SERVER_KEY", ""),
		MidtransClientKey:    getEnv("MIDTRANS_CLIENT_KEY", ""),
		MidtransProduction:   getEnv("MIDTRANS_IS_PRODUCTION", "false") == "true",
		CORSOrigin:           getEnv("CORS_ORIGIN", "http://localhost:3000"),
		AdminDefaultPassword: getEnv("ADMIN_DEFAULT_PASSWORD", "Admin@greensco1"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
