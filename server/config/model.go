package config

type Config struct {
	AppName       string              `mapstructure:"appname"`
	HttpServer    HttpServerConfig    `mapstructure:"http"`
	Observability ObservabilityConfig `mapstructure:"observability"`
	Realtime      RealtimeHubConfig   `mapstructure:"realtime-hub"`
	Websocket     WebsocketConfig     `mapstructure:"websocket"`
	Redis         RedisConfig         `mapstructure:"redis"`
	Postgres      PostgresConfig      `mapstructure:"postgres"`
	AuthToken     AuthTokenConfig     `mapstructure:"authentication"`
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
	PingInterval int `mapstructure:"ping-interval"`
	PongWait     int `mapstructure:"pong-wait"`
	WriteWait    int `mapstructure:"write-wait"`
	MessageSize  int `mapstructure:"message-size"`
	MaxMessages  int `mapstructure:"max-messages"`
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
