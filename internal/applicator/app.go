package applicator

import (
	"jumyste-app-backend/internal/ai"
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
	ResumeService  *service.ResumeService
	AuthHandler    *handler.AuthHandler
	UserHandler    *handler.UserHandler
	VacancyHandler *handler.VacancyHandler
	ResumeHandler  *handler.ResumeHandler
	AIClient       *ai.OpenAIClient
	//RedisClient *redis.Client
}

func NewApp() *App {
	logger.Log.Info("Initializing database...")
	database.InitDB()
	database.RunMigrations()

	//logger.Log.Info("Initializing Redis...")
	////redisClient := redisPkg.InitRedis()

	logger.Log.Info("Initializing AI client...")
	aiClient := ai.NewOpenAIClient()

	logger.Log.Info("Initializing repositories...")
	authRepo := repository.NewAuthRepository(database.DB)
	userRepo := repository.NewUserRepository(database.DB)
	vacancyRepo := repository.NewVacancyRepository(database.DB)

	logger.Log.Info("Initializing services...")
	authService := service.NewAuthService(authRepo)
	userService := service.NewUserService(userRepo)
	vacancyService := service.NewVacancyService(vacancyRepo)
	resumeService := service.NewResumeService(aiClient)

	logger.Log.Info("Initializing handlers...")
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	vacancyHandler := handler.NewVacancyHandler(vacancyService)
	resumeHandler := handler.NewResumeHandler(resumeService)

	logger.Log.Info("Application initialized successfully")

	return &App{
		AuthRepo:       authRepo,
		UserRepo:       userRepo,
		VacancyRepo:    vacancyRepo,
		AuthService:    authService,
		UserService:    userService,
		VacancyService: vacancyService,
		ResumeService:  resumeService,
		AuthHandler:    authHandler,
		UserHandler:    userHandler,
		VacancyHandler: vacancyHandler,
		ResumeHandler:  resumeHandler,
		AIClient:       aiClient,
		//RedisClient: redisClient,
	}
}
