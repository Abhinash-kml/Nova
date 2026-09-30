package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/abhinash-kml/nova/server/achievements"
	"github.com/abhinash-kml/nova/server/apiserver"
	"github.com/abhinash-kml/nova/server/auth"
	"github.com/abhinash-kml/nova/server/auth/providers"
	"github.com/abhinash-kml/nova/server/cache"
	"github.com/abhinash-kml/nova/server/channels"
	"github.com/abhinash-kml/nova/server/clans"
	"github.com/abhinash-kml/nova/server/comments"
	"github.com/abhinash-kml/nova/server/common"
	"github.com/abhinash-kml/nova/server/config"
	"github.com/abhinash-kml/nova/server/economy"
	"github.com/abhinash-kml/nova/server/infra"
	"github.com/abhinash-kml/nova/server/inventory"
	"github.com/abhinash-kml/nova/server/leaderboard"
	"github.com/abhinash-kml/nova/server/observability"
	"github.com/abhinash-kml/nova/server/posts"
	"github.com/abhinash-kml/nova/server/secretsmanager"
	"github.com/abhinash-kml/nova/server/social"
	"github.com/abhinash-kml/nova/server/stats"
	"github.com/abhinash-kml/nova/server/users"
	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/contrib/bridges/otelzap"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Variables for configs
var (
	configFileBase      string
	configFileType      string
	configFileSeparator = "."
	configFile          string
	migrateDirection    string
	migratetionSteps    int
)

func init() {
	// Parse config file name
	configFileBase = *flag.String("config", "config", "Specify config file to be used")

	// Parse config file type flag
	configFileType = *flag.String("configtype", "yaml", "Specify default config file type [ yaml | json ]")

	// Parse deployment environment flag
	deploymentEnvironment := flag.String("deployment", "local", "Specify default deployment environment [ local | staging | production ]")
	configFile = configFileBase + configFileSeparator + *deploymentEnvironment

	// Parse migration flag
	migrateDirection = *flag.String("migrate", "up", "Specify migrate direction [ up | down | none ]")

	// Parse migration steps
	migratetionSteps = *flag.Int("steps", 1, "Specify steps for migrations [ default = 1 ]")
}

func main() {
	// Listen for interrupt & kill signal
	globalCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// --- Parse flags ---
	flag.Parse()
	if flag.Parsed() {
		fmt.Println("Parsed all cli flags")
	}
	flag.Usage()

	fmt.Println("Migrate flag", migrateDirection)
	fmt.Println("Migrate steps", migratetionSteps)

	// 1. Load configs
	config.Initialize(configFile, configFileType, "./")
	if !config.Load() {
		log.Fatal("Failed to load configs....")
	}
	log.Printf("Loaded configs from file: %s", configFile+configFileSeparator+configFileType)
	config := config.GetInstance()

	// Open file for writing logs
	file, err := os.OpenFile("./logs/temp.log", os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		log.Fatal("Failed to open file for writing temporary logs: %w", err)
	}
	defer file.Close()

	// Setup opentelemetry
	shutdownFunc, err := observability.SetupOTelSDK(globalCtx, &config.Observability)
	if err != nil {

		log.Fatal("Failed to setup opentelemtry for observability. Error: %w", err)
	}
	// Call shutdown func for proper cleanup so we dont leak anything
	defer func() {
		err = errors.Join(shutdownFunc(globalCtx))
	}()

	// Setup logger
	fileSyncer := zapcore.AddSync(file)
	stdOutSyncer := os.Stdout
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	fileEncoder := zapcore.NewJSONEncoder(encoderConfig)
	logLevel := zap.NewAtomicLevelAt(zapcore.DebugLevel)
	fileCore := zapcore.NewCore(fileEncoder, fileSyncer, logLevel)
	stdOutCore := zapcore.NewCore(consoleEncoder, stdOutSyncer, logLevel)
	otelLogCore := otelzap.NewCore("nova", otelzap.WithLoggerProvider(observability.LoggerProvider()))
	teeCore := zapcore.NewTee(fileCore, stdOutCore, otelLogCore)
	logger := zap.New(teeCore).WithOptions(zap.AddCaller())
	defer logger.Sync()

	logger.Sugar().Infof("Server started at: %w", time.Now())

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
	postgresPool := infra.NewPostgressPgxPool(globalCtx, postgresDsn)

	// Ping redis to test connection
	result, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("Failed to ping connected redis. Error:", err)
	}
	fmt.Println("Redis ping result:", result)

	// Ping postgres to test connection
	err = postgresPool.Ping(context.Background())
	if err != nil {
		fmt.Println("Failed to ping connected postgres client. Error:", err)
	}

	// Setup redis instrumentation
	err = redisotel.InstrumentTracing(redisClient)
	if err != nil {
		panic("Failed to setup redis otel tracing")
	}

	// Create gin router engine
	globalRouter := gin.New()

	// Create and use CORS middleware
	corsMiddleware := cors.New(cors.Config{
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
	})
	globalRouter.Use(corsMiddleware)

	// Setup Otelgin metrics middleware
	globalRouter.Use(otelgin.Middleware("nova-server"))

	// Setup logging middleware
	globalRouter.Use(ginzap.Ginzap(logger, time.RFC3339, true))

	// Setup Error handling middleware
	globalRouter.Use(common.ErrorHandler)

	// Setup Auth middleware
	//globalRouter.Use(auth.Token())

	// Setup domains of interests
	// postgresDb := stdlib.OpenDBFromPool(postgresPool)
	// migrationManager, err := infra.NewMigrationManager(postgresDb, "./migrations/postgres", logger)
	// if err != nil {
	// 	logger.Fatal("Failed to create migration manager", zap.Error(err))
	// }

	// // Run migrations
	// err = migrationManager.MigrateWithLock(context.Background(), migrateDirection, migratetionSteps)
	// if err != nil {
	// 	logger.Fatal("Failed to run migrations with lock & rollback", zap.Error(err))
	// }

	// Setup oauth callback handler
	globalRouter.GET("/token", func(ctx *gin.Context) {
		code := ctx.Query("code")
		logger.Info("Auth Code", zap.String("Value", code))
		ctx.Status(200)
	})

	// Setup secrets manager
	smanager, err := secretsmanager.NewSecretsManager("./secrets/oauth.json")
	if err != nil {
		logger.Fatal("Failed to create secrets manager", zap.Error(err))
	}

	// Get google & facebook secrets
	fbSecret, found := smanager.Get("facebook")
	if !found {
		log.Fatal("Secrets for facebook oauth doesnt exist")
	}
	googleSecret, found := smanager.Get("google")
	if !found {
		log.Fatal("Secrets for google oauth doesnt exist")
	}

	// Setup auth engine
	engine := auth.NewSocialAuthEngine()
	auth.SetDefaultAuthEngine(engine)

	googleProvider, err := providers.NewGoogleProvider(
		context.Background(),
		logger,
		googleSecret.ClientID,
		googleSecret.ClientSecret,
		googleSecret.RedirectURI,
		googleSecret.Scope)
	if err != nil {
		logger.Fatal("Failed to create google oauth provider", zap.Error(err))
	}
	facebookProvider, err := providers.NewFacebookProvider(
		context.Background(),
		logger,
		fbSecret.ClientID,
		fbSecret.ClientSecret,
		fbSecret.RedirectURI,
		fbSecret.Scope)
	if err != nil {
		logger.Fatal("Failed to create facebook oauth provider", zap.Error(err))
	}
	authURL := googleProvider.AuthURL()
	logger.Info("Google Auth", zap.String("URL", authURL))
	authURL = facebookProvider.AuthURL()
	logger.Info("Facebook Auth", zap.String("URL", authURL))
	engine.Register(googleProvider)
	engine.Register(facebookProvider)
	engine.List()

	redisAuthStore := auth.NewRedisAuthStore(redisClient)
	auth.SetupJwtService(&config.AuthToken, redisAuthStore, logger)

	tieredCache := cache.NewTieredCache(redisClient, time.Second*30, time.Second*30, time.Second*10)

	// Setup users module
	{
		usersRepository := users.NewPostgresRepositoryFromPgxPool(postgresPool, config, logger, config.SeedsConfig.Users)
		if err = usersRepository.Seed(context.Background()); err != nil {
			logger.Error("Failed to seed users repository", zap.Error(err))
		}
		usersService := users.NewLocalUsersService(usersRepository, tieredCache, config, logger)
		usersController := users.NewController(usersService, config, logger)
		users.SetupRoutes(globalRouter, usersController)
	}

	// Setup posts module
	{
		postsRepository := posts.NewPostgresRepositoryFromPgxPool(postgresPool, config, logger, config.SeedsConfig.Posts)
		if err = postsRepository.Seed(context.Background()); err != nil {
			logger.Error("Failed to seed posts repository", zap.Error(err))
		}
		postsService := posts.NewLocalPostsService(postsRepository, tieredCache, config, logger)
		postsController := posts.NewController(postsService, config, logger)
		posts.SetupRoutes(globalRouter, postsController)
	}

	// Setup comments module
	{
		commentsRepository := comments.NewPostgresRepositoryFromPgxPool(postgresPool, config, logger, config.SeedsConfig.Comments)
		if err = commentsRepository.Seed(context.Background()); err != nil {
			logger.Error("Failed to seed comments repository", zap.Error(err))
		}
		commentsService := comments.NewLocalCommentsService(commentsRepository, config, redisClient, logger)
		commentsController := comments.NewController(commentsService, config, logger)
		comments.SetupRoutes(globalRouter, commentsController)
	}

	// Setup clans module
	{
		clansRepository := clans.NewPostgresRepositoryFromPgxPool(postgresPool, config, logger, config.SeedsConfig.Clans)
		if err = clansRepository.Seed(context.Background()); err != nil {
			logger.Error("Failed to seed clans repository", zap.Error(err))
		}
		clansService := clans.NewLocalClansService(clansRepository, config, redisClient, logger)
		clansController := clans.NewController(clansService, config, logger)
		clans.SetupRoutes(globalRouter, clansController)
	}

	// Setup channels module
	{
		channelsRepository := channels.NewPostgresRepositoryFromPgxPool(postgresPool, config, logger, config.SeedsConfig.Channels)
		if err = channelsRepository.Seed(context.Background()); err != nil {
			logger.Error("Failed to seed channels repository", zap.Error(err))
		}
		channelsService := channels.NewLocalChannelService(channelsRepository, config, logger)
		channelsController := channels.NewController(channelsService, config, logger)
		channels.SetupRoutes(globalRouter, channelsController)
	}

	// Setup stats module
	{
		statsRepository := stats.NewPostgresRepositoryFromPgxPool(postgresPool, config, logger, config.SeedsConfig.Stats)
		if err = statsRepository.Seed(context.Background()); err != nil {
			logger.Error("Failed to seed stats repository", zap.Error(err))
		}
		statsService := stats.NewService(statsRepository, config, logger)
		statsController := stats.NewController(statsService, config, logger)
		stats.SetupRoutes(globalRouter, statsController)
	}
	// Setup leaderboard module
	{
		leaderboardScoreRepo := leaderboard.NewRedisScoreRepository(config, logger, redisClient)
		leaderderboardMetaRepo := leaderboard.NewPostgresMetaRepository(postgresPool, config, logger, leaderboardScoreRepo)
		leaderboardService, _ := leaderboard.NewLocalService(globalCtx, config, logger, leaderderboardMetaRepo, leaderboardScoreRepo)
		leaderboardController := leaderboard.NewController(leaderboardService, config, logger)
		leaderboard.SetupRoutes(globalRouter, leaderboardController)
	}
	// Setup achievement module
	{
		// Achievement
		achievementRepo := achievements.NewPostgresRepository(postgresPool, config, logger)
		achievementService := achievements.NewAchievementService(achievementRepo, config, logger)
		achievementController := achievements.NewController(achievementService, config, logger)

		// Criteria
		criteriaRepo := achievements.NewCriteriaPostgresRepository(postgresPool, config, logger)
		criteriaService := achievements.NewCriteriaService(criteriaRepo, config, logger)
		criteriaController := achievements.NewCriteriaController(criteriaService, config, logger)

		// Progress
		progressRepo := achievements.NewProgressPostgresRepository(postgresPool, config, logger)
		progressService := achievements.NewLocalProgressServiceService(progressRepo, config, logger)
		progressController := achievements.NewProgressController(progressService, config, logger)

		// Completed
		completedRepo := achievements.NewPostgresCompletedAchievementRepository(postgresPool, config, logger)
		completedService := achievements.NewLocalCompletedAchievementService(completedRepo, config, logger)
		completedController := achievements.NewCompletedAchievementController(completedService, config, logger)

		// Setup routes
		achievements.SetupRoutes(globalRouter, achievementController, criteriaController, progressController, completedController)
	}
	// Setup inventory module
	{
		// Setup items
		itemsRepo := inventory.NewItemsPostgresRepositoryFromPgxPool(postgresPool, config, logger)
		itemsService := inventory.NewLocalItemsService(itemsRepo, config, logger)
		itemsController := inventory.NewItemsController(itemsService, config, logger)

		// Setup inventory
		inventoryRepo := inventory.NewPostgresInventoryRepositoryFromPgxPool(postgresPool, config, logger)
		inventoryService := inventory.NewLocalInventoryService(inventoryRepo, config, logger)
		inventoryController := inventory.NewInventoryController(inventoryService, config, logger)

		// Setup routes
		inventory.SetupRoutes(globalRouter, itemsController, inventoryController)
	}

	// Setup social module
	{
		socialRepository := social.NewPostgreRepository(postgresPool, config, logger)
		socialService := social.NewLocalSocialService(socialRepository, config, logger)
		socialController := social.NewController(socialService, config, logger)
		social.SetupRoutes(globalRouter, socialController)
	}

	//	Setup economy module
	{
		// Setup currency
		currencyRepo := economy.NewPostgresCurrencyRepository(postgresPool, config, logger)
		currencyService := economy.NewLocalCurrencyService(currencyRepo, config, logger)
		currencyController := economy.NBewCurrencyController(currencyService, config, logger)

		// Setup wallet
		walletRepo := economy.NewPostgresWalletRepository(postgresPool, config, logger)
		walletService := economy.NewLocalWalletService(walletRepo, config, logger)
		walletController := economy.NewWalletController(walletService, config, logger)

		economy.SetupRoutes(globalRouter, currencyController, walletController)
	}

	// Setup realtime module
	// realtimeBroker := realtime.NewRedisBroker(globalCtx, redisClient)
	// realtime.
	// realtime.NewHub(globalCtx, realtimeBroker, )

	// Create http api server & start it
	server := apiserver.New(globalCtx, config.HttpServer, globalRouter, logger)
	err = server.Start()
	if err != nil {
		logger.Error("Failed to start http api server", zap.Error(err))
	}

	// Block untill our signal is trigerred
	<-globalCtx.Done()
	shutdownFunc(context.Background())
	// Call stop() to immeaditely stop downstream services
	stop()
	postgresPool.Close()
	redisClient.Close()
}
