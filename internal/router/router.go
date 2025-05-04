package router

import (
	"github.com/gin-gonic/gin"
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
	invitationHandler *handler.InvitationHandler,
	jobApplicationHandler *handler.JobApplicationHandler,
	departmentHandler *handler.DepartmentsHandler,
) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	// Swagger documentation route
	//r.GET("api/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// --- Аутентификация ---
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", authHandler.RegisterUser)
		auth.POST("/register-hr", authHandler.RegisterHR)
		auth.POST("/login", authHandler.Login)
		auth.POST("/forgot-password", authHandler.RequestPasswordReset)
		auth.POST("/reset-password", authHandler.ResetPassword)
		auth.POST("/refresh", authHandler.RefreshToken)
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
		vacancyRoutes.GET("/company", vacancyHandler.GetVacancyByCompanyID)
		vacancyRoutes.POST("/", vacancyHandler.CreateVacancy)
		vacancyRoutes.PUT("/:id", vacancyHandler.UpdateVacancy)
		vacancyRoutes.DELETE("/:id", vacancyHandler.DeleteVacancy)
		vacancyRoutes.GET("/", vacancyHandler.GetAllVacancies)
		vacancyRoutes.GET("/my", vacancyHandler.GetMyVacancies)
		vacancyRoutes.GET("/search", vacancyHandler.SearchVacancies)
		vacancyRoutes.GET("/user/:id", vacancyHandler.GetVacancyByIDForUser)
		vacancyRoutes.GET("/hr/:id", vacancyHandler.GetVacancyByIDForHr)
		vacancyRoutes.GET("/feed/data", vacancyHandler.GetFeedData)
		vacancyRoutes.PUT("/status/:id", vacancyHandler.UpdateVacancyStatusHandler)
		vacancyRoutes.POST("/generate-description", vacancyHandler.GenerateVacancyDescription)
	}

	// --- Чаты ---
	chatRoutes := r.Group("/api/chats")
	chatRoutes.Use(authMiddleware.VerifyTokenMiddleware())
	{
		chatRoutes.POST("/", chatHandler.CreateChatHandler)
		chatRoutes.GET("/:chatID", chatHandler.GetChatByIDHandler)
		chatRoutes.GET("/", chatHandler.GetAllChatsHandler)
		chatRoutes.GET("/user", chatHandler.GetChatsByUserIDHandler)
	}

	// --- Сообщения ---
	messageRoutes := r.Group("/api/messages")
	messageRoutes.Use(authMiddleware.VerifyTokenMiddleware())
	{
		messageRoutes.POST("/", messageHandler.SendMessageHandler)
		messageRoutes.GET("/chat/:chatID", messageHandler.GetMessagesByChatIDHandler)
		messageRoutes.GET("/:messageID", messageHandler.GetMessageByIDHandler)
		messageRoutes.POST("/read", messageHandler.MarkAsRead)
	}

	// --- Резюме ---
	resume := r.Group("/api/resume")
	resume.Use(authMiddleware.VerifyTokenMiddleware())
	{
		resume.POST("/upload", resumeHandler.UploadResume)
		resume.POST("/manual", resumeHandler.CreateResume)
		resume.GET("/:user_id", resumeHandler.GetResumeByUserID)
		resume.DELETE("/", resumeHandler.DeleteResumeByUserID)
		resume.GET("/candidates", resumeHandler.FilterCandidates)
	}

	// --- Приглашения ---
	invitations := r.Group("/api/invitations")
	invitations.Use(authMiddleware.VerifyTokenMiddleware())
	{
		invitations.POST("/", invitationHandler.SendInvitationHandler)
	}

	// -- Отклики ---
	jobApp := r.Group("/api/jobs")
	jobApp.Use(authMiddleware.VerifyTokenMiddleware())
	{
		jobApp.POST("/apply/:vacancy_id", jobApplicationHandler.ApplyForJob)
		jobApp.GET("/:vacancy_id", jobApplicationHandler.GetJobApplicationsByVacancyID)
		jobApp.PUT("/:application_id/status/:status", jobApplicationHandler.UpdateJobApplicationStatus)
		jobApp.DELETE("/:application_id", jobApplicationHandler.DeleteJobApplication)
		jobApp.GET("/analytics", jobApplicationHandler.GetJobAppAnalytics)
	}

	departments := r.Group("/api/departments")
	departments.Use(authMiddleware.VerifyTokenMiddleware())
	{
		departments.GET("/all", departmentHandler.GetMyDepartments)
		departments.POST("/", departmentHandler.CreateDepartment)
		departments.GET("/:id", departmentHandler.GetDepartmentByID)
	}

	// --- WebSocket ---
	ws := r.Group("/api")
	ws.GET("/ws", wsHandler.HandleWebSocket)

	return r
}
