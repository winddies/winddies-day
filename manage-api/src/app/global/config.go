package global

import (
	"fmt"

	"github.com/spf13/viper"
)

var Conf Config

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type Config struct {
	Port     string      `mapstructure:"port"`
	Mode     Mode        `mapstructure:"mode"`
	Redis    RedisConfig `mapstructure:"redis"`
	LogLevel string      `mapstructure:"logLevel"`
}

type Mode string

const (
	DevMode  Mode = "dev"
	TestMode Mode = "test"
	ProdMode Mode = "prod"
)

func Init(configPath string) (err error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.SetConfigFile(configPath)
	err = viper.ReadInConfig() // Find and read the config file
	if err != nil {            // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	viper.Unmarshal(&Conf)
	return nil
}
