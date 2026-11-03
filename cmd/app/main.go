package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
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
	configs "github.com/abhinash-kml/nova/server/config"
	"github.com/abhinash-kml/nova/server/economy"
	"github.com/abhinash-kml/nova/server/events"
	"github.com/abhinash-kml/nova/server/infra"
	"github.com/abhinash-kml/nova/server/inventory"
	"github.com/abhinash-kml/nova/server/leaderboard"
	"github.com/abhinash-kml/nova/server/observability"
	"github.com/abhinash-kml/nova/server/posts"
	"github.com/abhinash-kml/nova/server/realtime"
	"github.com/abhinash-kml/nova/server/secretsmanager"
	"github.com/abhinash-kml/nova/server/social"
	"github.com/abhinash-kml/nova/server/stats"
	"github.com/abhinash-kml/nova/server/users"
	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
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
	configs.Initialize(configFile, configFileType, "./")
	if !configs.Load() {
		log.Fatal("Failed to load configs....")
	}
	log.Printf("Loaded configs from file: %s", configFile+configFileSeparator+configFileType)
	config := configs.GetInstance()
	configs.RegisterOnChangeHook(func(newConfig *configs.Config) {
		config = newConfig
	})

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
	var level zap.AtomicLevel
	level, err = zap.ParseAtomicLevel(config.Logging.Level)
	if err != nil {
		fmt.Println("Failed to parse logging level from configs:", err)
		fmt.Println("Setting default level to info")
		level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	logLevel := level
	fileCore := zapcore.NewCore(fileEncoder, fileSyncer, logLevel)
	stdOutCore := zapcore.NewCore(consoleEncoder, stdOutSyncer, logLevel)
	otelLogCore := otelzap.NewCore("nova", otelzap.WithLoggerProvider(observability.LoggerProvider()))
	teeCore := zapcore.NewTee(fileCore, stdOutCore, otelLogCore)
	logger := zap.New(teeCore).WithOptions(zap.AddCaller())
	zap.ReplaceGlobals(logger)
	defer logger.Sync()
	configs.RegisterOnChangeHook(func(c *configs.Config) {
		newLevel, err := zapcore.ParseLevel(c.Logging.Level)
		if err != nil {
			logger.Info("Failed to parse new incoming log level", zap.Error(err))
		} else {
			logger.Info("Logging level changed", zap.String("level", c.Logging.Level))
			level.SetLevel(newLevel)
		}
	})

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

	// Global events manager
	globalEventsManager := events.NewInMemoryEventsManager(globalCtx, 10)
	events.SetGlobalEventsManager(globalEventsManager)

	// Setup users module
	{
		usersRepository := users.NewPostgresRepositoryFromPgxPool(postgresPool, &config.UsersConfig, logger)
		usersRepository.Seed(context.Background())
		usersService := users.NewLocalUsersService(usersRepository, tieredCache, &config.UsersConfig, logger)
		usersController := users.NewController(usersService, config, logger)
		users.SetupRoutes(globalRouter, usersController)
	}

	// Setup posts module
	{
		postsRepository := posts.NewPostgresRepositoryFromPgxPool(postgresPool, &config.PostsConfig, logger)
		postsRepository.Seed(context.Background())
		postsService := posts.NewLocalPostsService(postsRepository, tieredCache, &config.PostsConfig, logger)
		postsController := posts.NewController(postsService, config, logger)
		posts.SetupRoutes(globalRouter, postsController)
	}

	// Setup comments module
	{
		commentsRepository := comments.NewPostgresRepositoryFromPgxPool(postgresPool, &config.CommentsConfig, logger)
		commentsRepository.Seed(context.Background())
		commentsService := comments.NewLocalCommentsService(commentsRepository, tieredCache, &config.CommentsConfig, logger)
		commentsController := comments.NewController(commentsService, config, logger)
		comments.SetupRoutes(globalRouter, commentsController)
	}

	// Setup clans module
	{
		clansRepository := clans.NewPostgresRepositoryFromPgxPool(postgresPool, &config.ClansConfig, logger)
		clansRepository.Seed(context.Background())
		clansService := clans.NewLocalClansService(clansRepository, &config.ClansConfig, tieredCache, logger)
		clansController := clans.NewController(clansService, config, logger)
		clans.SetupRoutes(globalRouter, clansController)
	}

	// Setup channels module
	{
		channelsRepository := channels.NewPostgresRepositoryFromPgxPool(postgresPool, &config.ChannelsConfig, logger)
		channelsRepository.Seed(context.Background())
		channelsService := channels.NewLocalChannelService(channelsRepository, config, logger)
		channelsController := channels.NewController(channelsService, &config.ChannelsConfig, logger)
		channels.SetupRoutes(globalRouter, channelsController)
	}

	// Setup stats module
	{
		statsRepository := stats.NewPostgresRepositoryFromPgxPool(postgresPool, &config.StatsConfig, logger)
		statsRepository.Seed(context.Background())
		statsService := stats.NewService(statsRepository, &config.StatsConfig, logger)
		statsService.Initialize()
		statsController := stats.NewController(statsService, &config.StatsConfig, logger)
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

		// Achievement evaluator
		evaluator := achievements.NewEvaluator(globalCtx, achievementService, "stats:update")
		evaluator.Initialize()
		evaluator.BuildMapping()
		evaluator.Run()

		// Setup routes
		achievements.SetupRoutes(globalRouter, achievementController)
	}
	// Setup inventory module
	{
		// Setup items
		itemsRepo := inventory.NewItemsPostgresRepositoryFromPgxPool(postgresPool, logger)
		itemsService := inventory.NewLocalItemsService(itemsRepo, config, logger)
		itemsController := inventory.NewItemsController(itemsService, logger)

		// Setup inventory
		inventoryRepo := inventory.NewPostgresInventoryRepositoryFromPgxPool(postgresPool, &config.InventoryConfig, logger)
		inventoryService := inventory.NewLocalInventoryService(inventoryRepo, config, logger)
		inventoryController := inventory.NewInventoryController(inventoryService, &config.InventoryConfig, logger)

		// Setup routes
		inventory.SetupRoutes(globalRouter, itemsController, inventoryController)
	}

	//	Setup economy module
	{
		// Setup currency
		currencyRepo := economy.NewPostgresCurrencyRepository(postgresPool, &config.CurrencyConfig, logger)
		currencyService := economy.NewLocalCurrencyService(currencyRepo, &config.CurrencyConfig, logger)
		currencyController := economy.NewCurrencyController(currencyService, &config.CurrencyConfig, logger)

		// Setup wallet
		walletRepo := economy.NewPostgresWalletRepository(postgresPool, &config.WalletConfig, logger)
		walletService := economy.NewLocalWalletService(walletRepo, &config.WalletConfig, logger)
		walletController := economy.NewWalletController(walletService, &config.WalletConfig, logger)

		economy.SetupRoutes(globalRouter, currencyController, walletController)
	}

	// Setup social module
	socialRepository := social.NewPostgreRepository(postgresPool, config, logger)
	socialService := social.NewLocalSocialService(socialRepository, config, logger)
	socialController := social.NewController(socialService, config, logger)
	social.SetupRoutes(globalRouter, socialController)

	persister := realtime.NewLocalMessagePersister(logger)
	realtime.SetGlobalPersistor(persister)
	ch := persister.Channel()
	listener := social.NewMessageListener(socialService, ch, logger)
	listener.Listen()

	// Setup realtime module
	// Setup broker
	realtimeBroker := realtime.NewRedisBroker(globalCtx, redisClient)
	// Setup session store
	realtimeSessionStore := realtime.NewInMemorySessionStore(100)
	// Setup hub
	realtimeHub := realtime.NewHub(globalCtx, realtimeBroker, realtimeSessionStore, config.Realtime, logger)
	realtimeHub.Initialize()
	realtimeHub.Run()
	err = realtime.SetGlobalHub(realtimeHub)
	if err != nil {
		panic(err)
	}

	// Presence manager
	presenceManager := realtime.NewPresenceManager(redisClient, realtimeHub.SendChannel())
	err = realtime.SetGlobalPresenceManager(presenceManager)
	if err != nil {
		panic(err)
	}

	// Setup ws upgrader and all
	globalRouter.GET("/ws", func(ctx *gin.Context) {
		upgrader := websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}

		userID, _ := uuid.Parse(ctx.GetHeader("userid"))

		conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, ctx.Request.Header)
		if err != nil {
			logger.Error("Failed to upgrade connection to websocket")
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "failed",
				"error":  err.Error(),
			})
			return
		}
		// defer conn.Close() <-- Dont do this here, this will close the connection on handler end

		rtClient := realtime.NewClient(&config.Websocket, userID, conn)
		realtimeSessionStore.Add(rtClient)
		realtimeHub.Register(rtClient)

		go rtClient.ReadIncoming()
		go rtClient.ProcessOutgoing()

		info := map[string]any{
			"user_id":       userID,
			"connected_at":  time.Now(),
			"ping_interval": config.Websocket.PingInterval, // Send this to client so clients dont need manual intervention
		}
		bytes, _ := json.Marshal(info)
		rtClient.Send(
			realtime.Envelope{
				Data: json.RawMessage(bytes),
			},
		)

		logger.Info("Realtime user connected", zap.String("user_id", userID.String()))
	})

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
