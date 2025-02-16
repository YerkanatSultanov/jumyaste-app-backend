package config

import (
	"log"

	"github.com/spf13/viper"
)

var AppConfig *Config

type Config struct {
	Server struct {
		Port int
	}

	Database struct {
		Host     string
		Port     int
		User     string
		Password string
		DBName   string
		SSLMode  string
	}

	JWT struct {
		Secret          string
		ExpirationHours int
	}

	GoogleOAuth struct {
		ClientID     string
		ClientSecret string
		RedirectURL  string
	}
}

func LoadConfig(path string) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Fatalf("Error unmarshalling config: %v", err)
	}

	log.Println("Configuration loaded successfully")
}
