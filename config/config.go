package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type DatabaseConfig struct {
	Host     string
	Name     string
	Username string
	Password string
	Port     string
	Option   string
}

func (c DatabaseConfig) URI() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?%s", c.Username, c.Password, c.Host, c.Port, c.Name, c.Option)
}

type RedisConfig struct {
	Host     string
	Password string
	DB       int
	Protocol int
}

type AppConfig struct {
	AudioDirectory string
	Database       DatabaseConfig
	Redis          RedisConfig
}

func Load() (AppConfig, error) {
	viper.SetConfigType("yaml")
	viper.SetConfigName("env.sample.yaml")
	viper.SetConfigName("env.yaml")
	viper.AddConfigPath(".")

	var c AppConfig

	err := viper.ReadInConfig()
	if err != nil {
		return c, nil
	}

	err = viper.Unmarshal(&c)
	if err != nil {
		return c, nil
	}

	return c, nil
}
