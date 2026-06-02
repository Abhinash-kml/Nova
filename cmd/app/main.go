package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/abhinash-kml/nova/server/config"
	"github.com/abhinash-kml/nova/server/infra"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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

	// Create & redis client instance
	redisClient := infra.NewRedis(redis.Options{
		Addr:     config.Redis.Address,
		DB:       config.Redis.Database,
		Username: config.Redis.Username,
		Password: config.Redis.Password,
	})

	// Create & connect postgres instance
	postgresDsn := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable",
		config.Postgres.Username,
		config.Postgres.Password,
		config.Postgres.Address,
		config.Postgres.Database)
	postgresClient := infra.NewPostgres(postgresDsn)

	// Ping redis to test connection
	result, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("Failed to ping connected redis", err)
	}
	fmt.Println("Redis ping result:", result)

	// Ping postgres to test connection
	err = postgresClient.Ping()
	if err != nil {
		fmt.Println("Failed to ping connected postgres client", err)
	}

	// Open file for writing logs
	file, err := os.OpenFile("./logs/temp.log", os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		log.Fatal("Failed to open file for writing temporary logs: %w", err)
	}
	defer file.Close()

	// Setup logger
	fileSyncer := zapcore.AddSync(file)
	stdOutSyncer := os.Stdout
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	fileEncoder := zapcore.NewJSONEncoder(encoderConfig)
	logLevel := zap.NewAtomicLevelAt(zapcore.InfoLevel)
	fileCore := zapcore.NewCore(fileEncoder, fileSyncer, logLevel)
	stdOutCore := zapcore.NewCore(consoleEncoder, stdOutSyncer, logLevel)
	teeCore := zapcore.NewTee(fileCore, stdOutCore)
	logger := zap.New(teeCore)
	defer logger.Sync()

	logger.Info("Test")
	logger.Error("Error", zap.Error(errors.New("Meow Meow")))

	// Block untill our signal is trigerred
	<-signalChan
	globalCtx.Done()

	// Gracefully shutdown all services by calling cancel() of global context
	cancel()
}
