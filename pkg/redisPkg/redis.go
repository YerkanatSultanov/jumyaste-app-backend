package redisPkg

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"jumyste-app-backend/pkg/logger"
	"os"
)

func InitRedis() *redis.Client {
	redisAddr := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	ctx := context.Background()
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		logger.Log.Error("Failed to connect to Redis:", err)
	}

	logger.Log.Info("Connected to Redis successfully", "address", redisAddr)
	return redisClient
}
