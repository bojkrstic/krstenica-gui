package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

type DBConfig struct {
	URL         string        `mapstructure:"url"`
	MaxIdleTime time.Duration `mapstructure:"maxidletime"`
	MaxLifetime time.Duration `mapstructure:"maxlifetime"`
	MaxOpenConn int           `mapstructure:"maxopenconn"`
	MaxIdleConn int           `mapstructure:"maxidleconn"`
}

// Config struct that helps parsing config.yaml
type Config struct {
	ENV string   `mapstructure:"env"`
	DB  DBConfig `mapstructure:"db"`

	HTTPPort       string          `mapstructure:"http_port"`
	JWTSecret      string          `mapstructure:"jwt_secret"`
	AdminJWTSecret string          `mapstructure:"admin_jwt_secret"`
	Host           string          `mapstructure:"host"`
	Migration      MigrationConfig `mapstructure:"migration"`
	Auth           AuthConfig      `mapstructure:"auth"`
}

type AuthConfig struct {
	Username      string `mapstructure:"username"`
	Password      string `mapstructure:"password"`
	SessionSecret string `mapstructure:"session_secret"`
}

func Load() (*Config, error) {
	var config Config

	viper.SetConfigName("config") // } config.yaml
	viper.SetConfigType("yaml")   // }
	viper.AddConfigPath(".")      // for local use
	viper.AddConfigPath("config")
	viper.AddConfigPath("./../../config") // for local use
	viper.AddConfigPath("./config")
	viper.SetEnvPrefix("krstenica") // set env vars prefix
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	viper.AutomaticEnv() // check ENV variables

	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
