package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"jumyste-app-backend/config"
	"jumyste-app-backend/pkg/logger"
	"log/slog"
)

var DB *sql.DB

func InitDB() {
	dbConfig := config.AppConfig.Database
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.DBName, dbConfig.SSLMode,
	)

	logger.Log.Info("Connecting to the AuthService database...",
		slog.String("host", dbConfig.Host),
		slog.String("dbname", dbConfig.DBName),
	)

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		logger.Log.Error("Failed to open AuthService database connection", slog.String("error", err.Error()))
		panic(err)
	}

	if err = DB.Ping(); err != nil {
		logger.Log.Error("Failed to ping AuthService database", slog.String("error", err.Error()))
		panic(err)
	}

	logger.Log.Info("Connected to AuthService database successfully")
}
