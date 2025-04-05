package redisPkg

import (
	"context"
	"github.com/redis/go-redis/v9"
	"jumyste-app-backend/pkg/logger"
	"os"
)

func InitRedis() *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),     // Пример: shortline.proxy.rlwy.net:20101
		Username: os.Getenv("REDIS_USERNAME"), // default
		Password: os.Getenv("REDIS_PASSWORD"), // IYcWcljTMbPGsEATRnPBzmRIAtSyKRlY
		DB:       0,
	})

	ctx := context.Background()
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		logger.Log.Error("Failed to connect to Redis", "address", os.Getenv("REDIS_ADDR"), "error", err)
		return nil
	}

	logger.Log.Info("Connected to Redis successfully", "address", os.Getenv("REDIS_ADDR"))
	return redisClient
}
