package config

import (
	"log"

	"github.com/spf13/viper"
)

func SetDefaults() {

}

func Initialize(name, filetype, path string) {
	viper.SetConfigName(name)
	viper.SetConfigType(filetype)
	viper.AddConfigPath(path)
}

func Load() bool {
	var configs Config
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Failed to read config. Error: %s", err.Error())
	}
	if err := viper.Unmarshal(&configs); err != nil {
		log.Fatalf("Failed to unmarshall config. Error: %s", err.Error())
	}

	instance = &configs

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
