package applicator

import (
	"jumyste-app-backend/internal/database"
	"jumyste-app-backend/internal/handler"
	"jumyste-app-backend/internal/repository"
	"jumyste-app-backend/internal/service"
	"jumyste-app-backend/pkg/logger"
	//"github.com/redis/go-redis/v9"
)

type App struct {
	AuthRepo       *repository.AuthRepository
	UserRepo       *repository.UserRepository
	VacancyRepo    *repository.VacancyRepository
	AuthService    *service.AuthService
	UserService    *service.UserService
	VacancyService *service.VacancyService
	AuthHandler    *handler.AuthHandler
	UserHandler    *handler.UserHandler
	VacancyHandler *handler.VacancyHandler
	//RedisClient *redis.Client
	ChatHandler    *handler.ChatHandler
	MessageHandler *handler.MessageHandler
}

func NewApp() *App {
	logger.Log.Info("Initializing database...")
	database.InitDB()
	database.RunMigrations()

	//logger.Log.Info("Initializing Redis...")
	////redisClient := redisPkg.InitRedis()

	logger.Log.Info("Initializing repositories...")
	authRepo := repository.NewAuthRepository(database.DB)
	userRepo := repository.NewUserRepository(database.DB)
	vacancyRepo := repository.NewVacancyRepository(database.DB)
	chatRepo := repository.NewChatRepository(database.DB)
	messageRepo := repository.NewMessageRepository(database.DB)

	logger.Log.Info("Initializing services...")
	authService := service.NewAuthService(authRepo)
	userService := service.NewUserService(userRepo)
	vacancyService := service.NewVacancyService(vacancyRepo)
	chatService := service.NewChatService(chatRepo)
	messageService := service.NewMessageService(messageRepo)

	logger.Log.Info("Initializing handlers...")
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	vacancyHandler := handler.NewVacancyHandler(vacancyService)
	chatHandler := handler.NewChatHandler(chatService)
	messageHandler := handler.NewMessageHandler(messageService)

	logger.Log.Info("Application initialized successfully")

	return &App{
		AuthRepo:       authRepo,
		UserRepo:       userRepo,
		AuthService:    authService,
		UserService:    userService,
		AuthHandler:    authHandler,
		UserHandler:    userHandler,
		VacancyHandler: vacancyHandler,
		ChatHandler:    chatHandler,
		MessageHandler: messageHandler,
		//RedisClient: redisClient,
	}
}
