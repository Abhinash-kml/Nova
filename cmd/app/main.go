package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/abhinash-kml/nova/server/config"
	"github.com/abhinash-kml/nova/server/infra"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Listen for interrupt & kill signal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// Global context for passing to all services
	globalCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Perform our task

	// 1. Load configs
	config.Initialize("config", "yaml", "./")
	if !config.Load() {
		log.Fatal("Failed to load configs....")
	}
	config := config.GetInstance()

	redisClient := infra.NewRedis(redis.Options{
		Addr:     config.Redis.Address,
		DB:       config.Redis.Database,
		Username: config.Redis.Username,
		Password: config.Redis.Password,
	})

	postgresDsn := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable",
		config.Postgres.Username,
		config.Postgres.Password,
		config.Postgres.Address,
		config.Postgres.Database)
	postgresClient := infra.NewPostgres(postgresDsn)

	result, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("Failed to ping connected redis", err)
	}
	fmt.Println("Redis ping result:", result)

	err = postgresClient.Ping()
	if err != nil {
		fmt.Println("Failed to ping connected postgres client", err)
	}

	// Block untill our signal is trigerred
	<-signalChan

	// Gracefully shutdown all services by calling cancel() of global context
	cancel()
}
