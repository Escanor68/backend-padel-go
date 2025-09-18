package config

import (
	"os"
	"strconv"
)

type Config struct {
	Database DatabaseConfig
	JWT      JWTConfig
	Server   ServerConfig
	MercadoPago MercadoPagoConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	Charset  string
}

type JWTConfig struct {
	SecretKey     string
	AccessTokenExpiry  int // en horas
	RefreshTokenExpiry int // en horas
}

type ServerConfig struct {
	Port string
	Mode string
}

type MercadoPagoConfig struct {
	AccessToken string
	WebhookSecret string
	FrontendURL string
	BackendURL  string
}

func Load() *Config {
	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "padel_db"),
			Charset:  getEnv("DB_CHARSET", "utf8mb4"),
		},
		JWT: JWTConfig{
			SecretKey:         getEnv("JWT_SECRET", "your-secret-key"),
			AccessTokenExpiry:  getEnvAsInt("JWT_ACCESS_EXPIRY", 24), // 24 horas
			RefreshTokenExpiry: getEnvAsInt("JWT_REFRESH_EXPIRY", 168), // 7 días
		},
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Mode: getEnv("GIN_MODE", "debug"),
		},
		MercadoPago: MercadoPagoConfig{
			AccessToken:  getEnv("MERCADOPAGO_ACCESS_TOKEN", ""),
			WebhookSecret: getEnv("MERCADOPAGO_WEBHOOK_SECRET", ""),
			FrontendURL:  getEnv("FRONTEND_URL", "http://localhost:3000"),
			BackendURL:   getEnv("BACKEND_URL", "http://localhost:8080"),
		},
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
