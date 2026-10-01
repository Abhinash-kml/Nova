package config

import (
	"fmt"
	"log"
	"sync"

	"github.com/fsnotify/fsnotify"
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
			log.Fatalf("Failed to read configs. Error: %s", err)
		}
		if err := viper.Unmarshal(&configs); err != nil {
			log.Fatalf("Failed to unmarshall configs. Error: %s", err)
		}
		viper.WatchConfig()

		instance = &configs
	})

	viper.OnConfigChange(func(in fsnotify.Event) {
		if err := viper.Unmarshal(instance); err != nil {
			fmt.Println("Failed to unmarshall updated config", err)
		} else {
			fmt.Println("Loaded new configs via hot-reload")
		}
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
