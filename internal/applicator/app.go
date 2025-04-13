package applicator

import (
	"github.com/redis/go-redis/v9"
	"jumyste-app-backend/internal/ai"
	"jumyste-app-backend/internal/database"
	"jumyste-app-backend/internal/handler"
	"jumyste-app-backend/internal/manager"
	"jumyste-app-backend/internal/middleware"
	"jumyste-app-backend/internal/repository"
	"jumyste-app-backend/internal/service"
	"jumyste-app-backend/pkg/logger"
	"jumyste-app-backend/pkg/redisPkg"
)

type App struct {
	AuthRepo          *repository.AuthRepository
	UserRepo          *repository.UserRepository
	VacancyRepo       *repository.VacancyRepository
	JobAppRepo        *repository.JobApplicationRepository // Добавлено
	AuthService       *service.AuthService
	UserService       *service.UserService
	VacancyService    *service.VacancyService
	ResumeService     *service.ResumeService
	InvitationService *service.InvitationService
	JobAppService     *service.JobApplicationService
	AuthHandler       *handler.AuthHandler
	UserHandler       *handler.UserHandler
	VacancyHandler    *handler.VacancyHandler
	JobAppHandler     *handler.JobApplicationHandler
	ResumeHandler     *handler.ResumeHandler
	InvitationHandler *handler.InvitationHandler
	AIClient          *ai.OpenAIClient
	ChatHandler       *handler.ChatHandler
	MessageHandler    *handler.MessageHandler
	WSManager         *manager.WebSocketManager
	WSHandler         *handler.WebSocketHandler
	RedisClient       *redis.Client
}

func NewApp(authMiddleware *middleware.AuthMiddleware) *App {
	logger.Log.Info("Initializing database...")
	database.InitDB()
	database.RunMigrations()

	logger.Log.Info("Initializing Redis client...")
	redisClient := redisPkg.InitRedis()

	logger.Log.Info("Initializing AI client...")

	aiClient := ai.NewOpenAIClient()

	logger.Log.Info("Initializing repositories...")
	authRepo := repository.NewAuthRepository(database.DB)
	userRepo := repository.NewUserRepository(database.DB)
	vacancyRepo := repository.NewVacancyRepository(database.DB)
	invitationRepo := repository.NewInvitationRepository(database.DB)
	hrRepo := repository.NewHrRepository(database.DB)
	chatRepo := repository.NewChatRepository(database.DB)
	messageRepo := repository.NewMessageRepository(database.DB)
	resumeRepo := repository.NewResumeRepository(database.DB)
	jobAppRepo := repository.NewJobApplicationRepository(database.DB)

	logger.Log.Info("Initializing services...")
	authService := service.NewAuthService(authRepo, redisClient, invitationRepo, hrRepo)
	userService := service.NewUserService(userRepo)
	vacancyService := service.NewVacancyService(vacancyRepo, aiClient)
	invitationService := service.NewInvitationService(invitationRepo)
	chatService := service.NewChatService(chatRepo)
	messageService := service.NewMessageService(messageRepo)
	resumeService := service.NewResumeService(aiClient, resumeRepo)
	jobAppService := service.NewJobApplicationService(jobAppRepo, resumeRepo, vacancyRepo, aiClient)

	logger.Log.Info("Initializing WebSocket manager...")
	wsManager := manager.NewWebSocketManager()
	go wsManager.Run()

	logger.Log.Info("Initializing handlers...")
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	vacancyHandler := handler.NewVacancyHandler(vacancyService)
	invitationHandler := handler.NewInvitationHandler(invitationService)
	chatHandler := handler.NewChatHandler(chatService)
	messageHandler := handler.NewMessageHandler(messageService, wsManager)
	resumeHandler := handler.NewResumeHandler(resumeService)
	wsHandler := handler.NewWebSocketHandler(wsManager, authMiddleware)
	jobAppHandler := handler.NewJobApplicationHandler(jobAppService, resumeService)

	logger.Log.Info("Application initialized successfully")

	return &App{
		AuthRepo:          authRepo,
		UserRepo:          userRepo,
		VacancyRepo:       vacancyRepo,
		JobAppRepo:        jobAppRepo,
		AuthService:       authService,
		UserService:       userService,
		VacancyService:    vacancyService,
		ResumeService:     resumeService,
		JobAppService:     jobAppService,
		InvitationService: invitationService,
		AuthHandler:       authHandler,
		UserHandler:       userHandler,
		VacancyHandler:    vacancyHandler,
		InvitationHandler: invitationHandler,
		ChatHandler:       chatHandler,
		MessageHandler:    messageHandler,
		ResumeHandler:     resumeHandler,
		JobAppHandler:     jobAppHandler,
		AIClient:          aiClient,
		WSManager:         wsManager,
		WSHandler:         wsHandler,
		RedisClient:       redisClient,
	}
}
