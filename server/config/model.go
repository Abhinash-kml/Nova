package config

import "time"

type Config struct {
	AppName         string              `mapstructure:"appname"`
	Logging         LoggingConfig       `mapstructure:"logging"`
	HttpServer      HttpServerConfig    `mapstructure:"http"`
	Observability   ObservabilityConfig `mapstructure:"observability"`
	Realtime        RealtimeHubConfig   `mapstructure:"realtime-hub"`
	Websocket       WebsocketConfig     `mapstructure:"websocket"`
	Redis           RedisConfig         `mapstructure:"redis"`
	Postgres        PostgresConfig      `mapstructure:"postgres"`
	AuthToken       AuthTokenConfig     `mapstructure:"authentication"`
	UsersConfig     UsersConfig         `mapstructure:"users"`
	PostsConfig     PostsConfig         `mapstructure:"posts"`
	CommentsConfig  CommentsConfig      `mapstructure:"comments"`
	ClansConfig     ClansConfig         `mapstructure:"clans"`
	ChannelsConfig  ChannelsConfig      `mapstructure:"channels"`
	StatsConfig     StatsConfig         `mapstructure:"stats"`
	ItemsConfig     ItemsConfig         `mapstructure:"playeritems"`
	InventoryConfig InventoryConfig     `mapstructure:"inventory"`
	CurrencyConfig  CurrencyConfig      `mapstructure:"currency"`
	WalletConfig    WalletConfig        `mapstructure:"wallets"`
}

type LoggingConfig struct {
	Level    string `mapstructure:"level"`
	FilePath string `mapstructure:"path"`
}

type HttpServerConfig struct {
	Address           string `mapstructure:"address"`
	ReadTimeout       int    `mapstructure:"readtimeout"`
	WriteTimeout      int    `mapstructure:"writetimeout"`
	IdleTimeout       int    `mapstructure:"idletimeout"`
	ReadHeaderTimeout int    `mapstructure:"readheadertimeout"`
	MaxHeaderBytes    int    `mapstructure:"maxheaderbytes"`
}

type ObservabilityConfig struct {
	Trace   OtelTraceConfig  `mapstructure:"trace"`
	Metrics OtelMetricConfig `mapstructure:"metric"`
	Logs    OtelLogConfig    `mapstructure:"log"`
}

type OtelTraceConfig struct {
	Endpoint string `mapstructure:"endpoint"`
}

type OtelMetricConfig struct {
	Endpoint string `mapstructure:"endpoint"`
}

type OtelLogConfig struct {
	Endpoint string `mapstructure:"endpoint"`
}

type RealtimeHubConfig struct {
	Goroutine RealtimeHubGoroutineConfig `mapstructure:"goroutine"`
}

type RealtimeHubGoroutineConfig struct {
	MaxMainGoroutine   int `mapstructure:"main"`
	MaxRouterGoroutine int `mapstructure:"router"`
	MaxBrokerGoroutine int `mapstructure:"broker"`
}

type WebsocketConfig struct {
	PingInterval time.Duration `mapstructure:"ping-interval"`
	PongWait     time.Duration `mapstructure:"pong-wait"`
	WriteWait    time.Duration `mapstructure:"write-wait"`
	MessageSize  int64         `mapstructure:"message-size"`
	MaxMessages  int64         `mapstructure:"max-messages"`
}

type RedisConfig struct {
	Address  string `mapstructure:"address"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database int    `mapstructure:"database"`
}

type PostgresConfig struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Address  string `mapstructure:"address"`
	Database string `mapstructure:"database"`
}

type AuthTokenConfig struct {
	AccessToken  TokenConfig `mapstructure:"access"`
	RefreshToken TokenConfig `mapstructure:"refresh"`
}

type TokenConfig struct {
	Secret    string `mapstructure:"secret"`
	ExpiresIn int64  `mapstructure:"expires_in"`
	Issuer    string `mapstructure:"issuer"`
	Audience  string `mapstructure:"audience"`
}

type UsersConfig struct {
	KeyPrefix     string        `mapstructure:"keyPrefix"`
	SeedFile      string        `mapstructure:"seed"`
	CacheDuration time.Duration `mapstructure:"cacheduration"`
}

type PostsConfig struct {
	KeyPrefix     string        `mapstructure:"keyPrefix"`
	SeedFile      string        `mapstructure:"seed"`
	CacheDuration time.Duration `mapstructure:"cacheduration"`
}

type CommentsConfig struct {
	KeyPrefix     string        `mapstructure:"keyPrefix"`
	SeedFile      string        `mapstructure:"seed"`
	CacheDuration time.Duration `mapstructure:"cacheduration"`
}

type ClansConfig struct {
	KeyPrefix           string        `mapstructure:"keyPrefix"`
	ClanSeedFile        string        `mapstructure:"seed"`
	ClanMembersSeedFile string        `mapstructure:"members-seed"`
	CacheDuration       time.Duration `mapstructure:"cacheduration"`
}

type ChannelsConfig struct {
	KeyPrefix     string        `mapstructure:"keyPrefix"`
	SeedFile      string        `mapstructure:"seed"`
	CacheDuration time.Duration `mapstructure:"cacheduration"`
}

type StatsConfig struct {
	KeyPrefix           string        `mapstructure:"keyPrefix"`
	SeedFile            string        `mapstructure:"seed"`
	ProgressionSeedFile string        `mapstructure:"progression-seed"`
	CacheDuration       time.Duration `mapstructure:"cacheduration"`
}

type InventoryConfig struct {
	KeyPrefix     string        `mapstructure:"keyPrefix"`
	SeedFile      string        `mapstructure:"seed"`
	CacheDuration time.Duration `mapstructure:"cacheduration"`
}

type ItemsConfig struct {
	KeyPrefix     string        `mapstructure:"keyPrefix"`
	SeedFile      string        `mapstructure:"seed"`
	CacheDuration time.Duration `mapstructure:"cacheduration"`
}

type CurrencyConfig struct {
	KeyPrefix     string        `mapstructure:"keyPrefix"`
	SeedFile      string        `mapstructure:"seed"`
	CacheDuration time.Duration `mapstructure:"cacheduration"`
}

type WalletConfig struct {
	KeyPrefix     string        `mapstructure:"keyPrefix"`
	SeedFile      string        `mapstructure:"seed"`
	CacheDuration time.Duration `mapstructure:"cacheduration"`
}

type AchievementsConfig struct {
	KeyPrefix     string        `mapstructure:"keyPrefix"`
	SeedFile      string        `mapstructure:"seed"`
	CacheDuration time.Duration `mapstructure:"cacheduration"`
}
