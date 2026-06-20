package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/abhinash-kml/nova/server/apiserver"
	"github.com/abhinash-kml/nova/server/channels"
	"github.com/abhinash-kml/nova/server/clans"
	"github.com/abhinash-kml/nova/server/comments"
	"github.com/abhinash-kml/nova/server/config"
	"github.com/abhinash-kml/nova/server/infra"
	"github.com/abhinash-kml/nova/server/observability"
	"github.com/abhinash-kml/nova/server/posts"
	"github.com/abhinash-kml/nova/server/users"
	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// Listen for interrupt & kill signal
	globalCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

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

	err := redisotel.InstrumentTracing(redisClient)
	if err != nil {
		panic("Failed to setup redis otel tracing")
	}
	err = redisotel.InstrumentMetrics(redisClient)
	if err != nil {
		panic("Failed to setup redis otel metric")
	}

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

	logger.Sugar().Infof("Current Time: %w", time.Now())

	// Setup opentelemetry
	shutdownFunc, err := observability.SetupOTelSDK(globalCtx)
	if err != nil {
		logger.Fatal("Failed to setup opentelemtry for observability", zap.Error(err))
	}
	// Call shutdown func for proper cleanup so we dont leak anything
	defer func() {
		err = errors.Join(shutdownFunc(globalCtx))
	}()

	// Create gin router engine
	globalRouter := gin.New()

	// Setup cors middleware
	globalRouter.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		// AllowOrigins: []string{""}, // Only in production
		AllowMethods:     []string{"GET", "POST", "PATCH", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowWebSockets:  true,
		MaxAge:           time.Hour * 12,
		// AllowOriginFunc: func(origin string) bool {
		// return origin == ""
		// },
	}))

	// Setup logging middleware
	globalRouter.Use(ginzap.Ginzap(logger, time.RFC3339, true))

	// Setup Auth middleware
	//globalRouter.Use(auth.Token())

	// Setup domains of interests

	// Setup users module
	usersTracer := otel.Tracer("users-domain")
	usersRepository := users.NewInMemoryUsersRepository(logger, usersTracer)
	usersRepository.Seed()
	usersService := users.NewLocalUsersService(usersRepository, redisClient, logger, usersTracer)
	usersController := users.NewController(usersService, logger, usersTracer)
	users.SetupRoutes(globalRouter, usersController)

	// Setup posts module
	postsTracer := otel.Tracer("posts-domain")
	postsRepository := posts.NewInMemoryPostsRepository(logger, postsTracer)
	postsRepository.Seed()
	postsService := posts.NewLocalPostsService(postsRepository, redisClient, logger, postsTracer)
	postsController := posts.NewController(postsService, logger, postsTracer)
	posts.SetupRoutes(globalRouter, postsController)

	// Setup comments module
	commentsTracer := otel.Tracer("comments-tracer")
	commentsRepository := comments.NewInMemoryCommentsRepository(logger, commentsTracer)
	commentsRepository.Seed()
	commentsService := comments.NewLocalCommentsService(commentsRepository, logger, commentsTracer)
	commentsController := comments.NewController(commentsService, logger, commentsTracer)
	comments.SetupRoutes(globalRouter, commentsController)

	// Setup clans module
	clansTracer := otel.Tracer("clans-tracer")
	clansRepository := clans.NewInMemoryClanRepository(logger, clansTracer)
	clansRepository.Seed()
	clansService := clans.NewLocalClansService(clansRepository, logger, clansTracer)
	clansController := clans.NewController(clansService, logger, clansTracer)
	clans.SetupRoutes(globalRouter, clansController)

	// Setup channels module
	channelsTracer := otel.Tracer("channels-tracer")
	channelsRepository := channels.NewInMemoryChannelsRepository(logger, channelsTracer)
	channelsRepository.Seed()
	channelsService := channels.NewLocalChannelService(channelsRepository, logger, channelsTracer)
	channelsController := channels.NewController(channelsService, logger, channelsTracer)
	channels.SetupRoutes(globalRouter, channelsController)

	// Create http api server & start it
	server := apiserver.New(globalCtx, config.HttpServer, globalRouter, logger)
	err = server.Start()
	if err != nil {
		logger.Error("Failed to start http api server", zap.Error(err))
	}

	// Block untill our signal is trigerred
	<-globalCtx.Done()

	// Call stop() to immeaditely stop downstream services
	stop()
}
