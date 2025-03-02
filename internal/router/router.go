package router

import (
	"github.com/gin-gonic/gin"
	"jumyste-app-backend/internal/handler"
	"jumyste-app-backend/internal/middleware"
)

func SetupRouter(authHandler *handler.AuthHandler, userHandler *handler.UserHandler, authMiddleware *middleware.AuthMiddleware) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	auth := r.Group("/api/auth")
	{
		auth.POST("/register", authHandler.Register)
		//auth.POST("/verify-code", authHandler.VerifyCodeAndRegister)
		auth.POST("/login", authHandler.Login)
		auth.POST("/forgot-password", authHandler.RequestPasswordReset)
		auth.POST("/reset-password", authHandler.ResetPassword)
	}

	protected := r.Group("/api/users")
	protected.Use(authMiddleware.VerifyTokenMiddleware())

	{
		protected.GET("/me", userHandler.GetUser)
		protected.PATCH("/me", userHandler.UpdateUser)
		// users.DELETE("/me", userHandler.DeleteUser)
	}

	return r
}
