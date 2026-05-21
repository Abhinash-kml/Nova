package config

import "sync"

var (
	instance *Config
	once     sync.Once
)

type Config struct {
	AppName             string              `mapstructure:"appname"`
	HttpServerConfig    HttpServerConfig    `mapstructure:"http"`
	ObservabilityConfig ObservabilityConfig `mapstructure:"observability"`
	RealtimeConfig      RealtimeHubConfig   `mapstructure:"realtime-hub"`
	WebsocketConfig     WebsocketConfig     `mapstructure:"websocket"`
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
