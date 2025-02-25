package applicator

import (
	"jumyste-app-backend/internal/database"
	"jumyste-app-backend/internal/handler"
	"jumyste-app-backend/internal/repository"
	"jumyste-app-backend/internal/service"
	"jumyste-app-backend/pkg/logger"
	"jumyste-app-backend/pkg/redisPkg"

	"github.com/redis/go-redis/v9"
)

type App struct {
	AuthRepo    *repository.AuthRepository
	UserRepo    *repository.UserRepository
	AuthService *service.AuthService
	UserService *service.UserService
	AuthHandler *handler.AuthHandler
	UserHandler *handler.UserHandler
	RedisClient *redis.Client
}

func NewApp() *App {
	logger.Log.Info("Initializing database...")
	database.InitDB()
	database.RunMigrations()

	logger.Log.Info("Initializing Redis...")
	redisClient := redisPkg.InitRedis()

	logger.Log.Info("Initializing repositories...")
	authRepo := repository.NewAuthRepository(database.DB, redisClient)
	userRepo := repository.NewUserRepository(database.DB)

	logger.Log.Info("Initializing services...")
	authService := service.NewAuthService(authRepo, redisClient)
	userService := service.NewUserService(userRepo)

	logger.Log.Info("Initializing handlers...")
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)

	logger.Log.Info("Application initialized successfully")

	return &App{
		AuthRepo:    authRepo,
		UserRepo:    userRepo,
		AuthService: authService,
		UserService: userService,
		AuthHandler: authHandler,
		UserHandler: userHandler,
		RedisClient: redisClient,
	}
}
