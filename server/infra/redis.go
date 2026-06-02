package infra

import (
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func NewRedis(config redis.Options) *redis.Client {
	client := redis.NewClient(&config)
	if client == nil {
		zap.L().Fatal("Failed to connect to redis")
	}

	return client
}
