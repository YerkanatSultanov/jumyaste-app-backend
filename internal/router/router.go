package router

import (
	"github.com/gin-gonic/gin"
	"jumyste-app-backend/internal/handler"
	"jumyste-app-backend/internal/middleware"
)

func SetupRouter(
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	vacancyHandler *handler.VacancyHandler,
	chatHandler *handler.ChatHandler,
	messageHandler *handler.MessageHandler,

	authMiddleware *middleware.AuthMiddleware,
) *gin.Engine {
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

	vacancyRoutes := r.Group("/api/vacancies")
	vacancyRoutes.Use(middleware.RequireRole(2))
	{
		vacancyRoutes.POST("/", vacancyHandler.CreateVacancy)
		vacancyRoutes.PUT("/:id", vacancyHandler.UpdateVacancy)
		vacancyRoutes.DELETE("/:id", vacancyHandler.DeleteVacancy)
		vacancyRoutes.GET("/", vacancyHandler.GetAllVacancies)
		vacancyRoutes.GET("/my", vacancyHandler.GetMyVacancies)
		vacancyRoutes.GET("/search", vacancyHandler.SearchVacancies)
	}

	chatRoutes := r.Group("/api/chats")
	chatRoutes.Use(middleware.AuthMiddlewares())
	{
		chatRoutes.POST("/", chatHandler.CreateChatHandler)
		chatRoutes.GET("/:chatID", chatHandler.GetChatByIDHandler)
		chatRoutes.GET("/", chatHandler.GetAllChatsHandler)
	}

	messageRoutes := r.Group("/api/messages")
	messageRoutes.Use(middleware.AuthMiddlewares())
	{
		messageRoutes.POST("/", messageHandler.SendMessageHandler)
		messageRoutes.GET("/chat/:chatID", messageHandler.GetMessagesByChatIDHandler)
		messageRoutes.GET("/:messageID", messageHandler.GetMessageByIDHandler)
  }
	resume := r.Group("/api/resume")
	resume.Use(authMiddleware.VerifyTokenMiddleware())
	{
		resume.POST("/upload", resumeHandler.UploadResume)
	}
	return r
}
