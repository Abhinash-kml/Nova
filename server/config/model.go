package config

import "sync"

var (
	instance *Config
	once     sync.Once
)

type Config struct {
	AppName          string           `mapstructure:"appname"`
	HttpServerConfig HttpServerConfig `mapstructure:"http"`
}

type HttpServerConfig struct {
	Address           string `mapstructure:"address"`
	ReadTimeout       int    `mapstructure:"readtimeout"`
	WriteTimeout      int    `mapstructure:"writetimeout"`
	IdleTimeout       int    `mapstructure:"idletimeout"`
	ReadHeaderTimeout int    `mapstructure:"readheadertimeout"`
	MaxHeaderBytes    int    `mapstructure:"maxheaderbytes"`
}
