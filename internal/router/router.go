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
		auth.POST("/login", authHandler.Login)
		auth.POST("/forgot-password", func(c *gin.Context) {
			authHandler.ForgotPasswordHandler(c.Writer, c.Request)
		})

		auth.POST("/reset-password", func(c *gin.Context) {
			authHandler.ResetPasswordHandler(c.Writer, c.Request)
		})
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
