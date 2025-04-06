// @title Jumyste App API
// @version 1.0
// @description API for Jumyste application
// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey	BearerAuth
// @type						apiKey
// @name						Authorization
// @in							header
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

	config.LoadConfig()

	logger.Log.Info("Starting application...")

	auth := middleware.NewAuthMiddleware(config.AppConfig)
	app := applicator.NewApp(auth)

	r := router.SetupRouter(
		app.AuthHandler,
		app.UserHandler,
		app.VacancyHandler,
		app.ChatHandler,
		app.MessageHandler,
		app.ResumeHandler,
		auth,
		app.WSHandler,
		app.InvitationHandler,
		app.JobAppHandler,
	)

	serverPort := config.AppConfig.Server.Port
	addr := fmt.Sprintf(":%s", serverPort)
	logger.Log.Info("Starting server", "port", serverPort)

	if err := r.Run(addr); err != nil {
		logger.Log.Error("Failed to start server", "error", err.Error())
	}
}
