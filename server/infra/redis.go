package infra

import (
	"log"

	"github.com/redis/go-redis/v9"
)

func NewRedis(config redis.Options) *redis.Client {
	client := redis.NewClient(&config)
	if client == nil {
		log.Fatal("Failed to connect to redis")
	}

	return client
}
