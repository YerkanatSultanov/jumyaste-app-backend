// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and your JWT token.
package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/swaggo/files"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"jumyste-app-backend/config"
	"jumyste-app-backend/docs"
	"jumyste-app-backend/internal/applicator"
	"jumyste-app-backend/internal/middleware"
	"jumyste-app-backend/internal/router"
	"jumyste-app-backend/pkg/logger"
	"os"
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
		app.DepartmentHandler,
	)

	serverPort := config.AppConfig.Server.Port
	addr := fmt.Sprintf(":%s", serverPort)
	logger.Log.Info("Starting server", "port", serverPort)
	setupSwagger(r)
	if err := r.Run(addr); err != nil {
		logger.Log.Error("Failed to start server", "error", err.Error())
	}

}

func setupSwagger(r *gin.Engine) {
	docs.SwaggerInfo.Title = "Jumyste App API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Description = "API for Jumyste application"

	env := os.Getenv("APP_ENV")
	if env == "production" {
		docs.SwaggerInfo.Host = "jumyaste-app-backend-production.up.railway.app"
		docs.SwaggerInfo.BasePath = "/api"
		docs.SwaggerInfo.Schemes = []string{"https"}
	} else {
		docs.SwaggerInfo.Host = "localhost:8080"
		docs.SwaggerInfo.BasePath = "/api"
		docs.SwaggerInfo.Schemes = []string{"http"}
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
