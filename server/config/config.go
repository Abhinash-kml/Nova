package config

import (
	"sync"

	"github.com/spf13/viper"
	"go.uber.org/zap"
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
			zap.L().Fatal("Failed to configs", zap.Error(err))
		}
		if err := viper.Unmarshal(&configs); err != nil {
			zap.L().Fatal("Faild to unmarshall configs", zap.Error(err))
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
