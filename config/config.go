package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	SMTP     SMTPConfig
	AI       AIConfig
	AppEnv   AppEnv
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type JWTConfig struct {
	Secret          string
	ExpirationHours int
}

type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Sender   string
}

type AIConfig struct {
	APIKey string
}

type AppEnv struct {
	AppEnv string
}

var AppConfig Config

func LoadConfig() {
	if err := godotenv.Load("./.env"); err != nil {
		log.Println("No .env file found, loading from environment")
	}

	AppConfig = Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "trolley.proxy.rlwy.net"),
			Port:     getEnv("DB_PORT", "35282"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "SXiJlyYzsIMXgeoMLkaYdnfidyuaitPN"),
			DBName:   getEnv("DB_NAME", "railway"),
			SSLMode:  getEnv("DB_SSLMODE", "require"),
		},
		JWT: JWTConfig{
			Secret:          getEnv("JWT_SECRET", "secretkey"),
			ExpirationHours: getEnvInt("JWT_EXPIRATION_HOURS", 1),
		},
		SMTP: SMTPConfig{
			Host:     getEnv("SMTP_HOST", "localhost"),
			Port:     getEnv("SMTP_PORT", "1025"),
			Username: getEnv("SMTP_USERNAME", ""),
			Password: getEnv("SMTP_PASSWORD", ""),
			Sender:   getEnv("SMTP_SENDER", "noreply@jumyste-app.local"),
		},
		AI: AIConfig{
			APIKey: getEnv("OPENAI_API_KE", ""),
		},
		AppEnv: AppEnv{
			AppEnv: getEnv("APP_ENV", "development"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
