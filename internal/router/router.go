package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	_ "jumyste-app-backend/docs"
	"jumyste-app-backend/internal/handler"
	"jumyste-app-backend/internal/middleware"
)

func SetupRouter(
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	vacancyHandler *handler.VacancyHandler,
	chatHandler *handler.ChatHandler,
	messageHandler *handler.MessageHandler,
	resumeHandler *handler.ResumeHandler,
	authMiddleware *middleware.AuthMiddleware,
	wsHandler *handler.WebSocketHandler,
) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	r.GET("api/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// --- Аутентификация ---
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", authHandler.Register)
		//auth.POST("/verify-code", authHandler.VerifyCodeAndRegister)
		auth.POST("/login", authHandler.Login)
		auth.POST("/forgot-password", authHandler.RequestPasswordReset)
		auth.POST("/reset-password", authHandler.ResetPassword)
	}

	// --- Пользователи ---
	protected := r.Group("/api/users")
	protected.Use(authMiddleware.VerifyTokenMiddleware())

	{
		protected.GET("/me", userHandler.GetUser)
		protected.PATCH("/me", userHandler.UpdateUser)
		// users.DELETE("/me", userHandler.DeleteUser)
	}

	// --- Вакансии (только для роли 2) ---
	vacancyRoutes := r.Group("/api/vacancies")
	vacancyRoutes.Use(authMiddleware.VerifyTokenMiddleware())
	vacancyRoutes.Use(middleware.RequireRole(2))
	{
		vacancyRoutes.POST("/", vacancyHandler.CreateVacancy)
		vacancyRoutes.PUT("/:id", vacancyHandler.UpdateVacancy)
		vacancyRoutes.DELETE("/:id", vacancyHandler.DeleteVacancy)
		vacancyRoutes.GET("/", vacancyHandler.GetAllVacancies)
		vacancyRoutes.GET("/my", vacancyHandler.GetMyVacancies)
		vacancyRoutes.GET("/search", vacancyHandler.SearchVacancies)
	}

	// --- Чаты ---
	chatRoutes := r.Group("/api/chats")
	chatRoutes.Use(authMiddleware.VerifyTokenMiddleware()) // Было неправильно
	{
		chatRoutes.POST("/", chatHandler.CreateChatHandler)
		chatRoutes.GET("/:chatID", chatHandler.GetChatByIDHandler)
		chatRoutes.GET("/", chatHandler.GetAllChatsHandler)
	}

	// --- Сообщения ---
	messageRoutes := r.Group("/api/messages")
	messageRoutes.Use(authMiddleware.VerifyTokenMiddleware()) // Было неправильно
	{
		messageRoutes.POST("/", messageHandler.SendMessageHandler)
		messageRoutes.GET("/chat/:chatID", messageHandler.GetMessagesByChatIDHandler)
		messageRoutes.GET("/:messageID", messageHandler.GetMessageByIDHandler)
		messageRoutes.POST("/read", messageHandler.MarkAsRead) // Фикс пути
	}

	// --- Резюме ---
	resume := r.Group("/api/resume")
	resume.Use(authMiddleware.VerifyTokenMiddleware())
	{
		resume.POST("/upload", resumeHandler.UploadResume)
	}

	// --- WebSocket ---
	ws := r.Group("/api")
	ws.GET("/ws", wsHandler.HandleWebSocket)

	return r
}
