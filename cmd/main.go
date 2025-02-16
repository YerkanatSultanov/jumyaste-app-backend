package main

import (
	"fmt"
	"jumyste-app-backend/config"
	"jumyste-app-backend/internal/applicator"
	"jumyste-app-backend/internal/middleware"
	"jumyste-app-backend/internal/router"
	"jumyste-app-backend/pkg/logger"
)

func main() {
	logger.InitLogger()

	config.LoadConfig("config/config.yaml")

	logger.Log.Info("Starting application...")
	app := applicator.NewApp()

	auth := middleware.NewAuthMiddleware(config.AppConfig)

	r := router.SetupRouter(app.AuthHandler, app.UserHandler, auth)

	serverPort := config.AppConfig.Server.Port
	addr := fmt.Sprintf(":%d", serverPort)
	logger.Log.Info("Starting server", "port", serverPort)

	if err := r.Run(addr); err != nil {
		logger.Log.Error("Failed to start server", "error", err.Error())
	}
}
