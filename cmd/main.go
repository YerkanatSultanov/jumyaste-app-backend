package main

import (
	"fmt"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	app := applicator.NewApp()

	auth := middleware.NewAuthMiddleware(config.AppConfig)

	r := router.SetupRouter(
		app.AuthHandler,
		app.UserHandler,
		app.VacancyHandler,
		app.ChatHandler,
		app.MessageHandler,
		auth,
	)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	serverPort := config.AppConfig.Server.Port
	addr := fmt.Sprintf(":%s", serverPort)
	logger.Log.Info("Starting server", "port", serverPort)

	if err := r.Run(addr); err != nil {
		logger.Log.Error("Failed to start server", "error", err.Error())
	}
}
