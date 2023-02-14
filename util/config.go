package util

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DBDriver string `mapstructure:"DB_Driver"`
	DBsource string `mapstructure:"DB_SOURCE"`
	Port     string `mapstructure:"PORT"`
}

var C *Config

func LoadConfig(path string) *Config {
	if C == nil {
		viper.AddConfigPath(path)
		viper.SetConfigName("app")
		viper.SetConfigType("env")

		viper.AutomaticEnv()

		if err := viper.ReadInConfig(); err != nil {
			log.Panicln("Err! cannot load config. Err : ", err)
		}

		if err := viper.Unmarshal(&C); err != nil {
			log.Panicln("Err! cannot unmarshal config. Err : ", err)
		}
	}

	return C
}
