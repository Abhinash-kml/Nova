package config

import (
	"log"
	"sync"

	"github.com/spf13/viper"
)

var (
	instance *Config
	once     sync.Once
)

func SetDefaults() {

}

func Initialize(name, filetype, path string) {
	viper.SetConfigName(name)
	viper.SetConfigType(filetype)
	viper.AddConfigPath(path)
}

func Load() bool {
	once.Do(func() {
		var configs Config
		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("Failed to read config. Error: %s", err.Error())
		}
		if err := viper.Unmarshal(&configs); err != nil {
			log.Fatalf("Failed to unmarshall config. Error: %s", err.Error())
		}

		instance = &configs
	})

	return true
}

func Reload() bool {
	return true
}

func GetInstance() *Config {
	if instance == nil {
		Load()
	}
	return instance
}
