package config

import (
	"net"
	neturl "net/url"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type DBConfig struct {
	URL         string        `mapstructure:"url"`
	LocalURL    string        `mapstructure:"local_url"`
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

	config.applyDefaults()

	return &config, nil
}

func (c *Config) applyDefaults() {
	c.DB.URL = strings.TrimSpace(c.DB.URL)
	c.DB.LocalURL = strings.TrimSpace(c.DB.LocalURL)

	if c.DB.URL == "" && c.DB.LocalURL != "" {
		c.DB.URL = c.DB.LocalURL
		return
	}

	if c.shouldUseLocalDBURL() && c.DB.LocalURL != "" {
		c.DB.URL = c.DB.LocalURL
	}
}

func (c *Config) shouldUseLocalDBURL() bool {
	env := strings.ToLower(strings.TrimSpace(c.ENV))
	switch env {
	case "local", "dev", "development":
		return true
	}

	host := extractDBHost(c.DB.URL)
	if host == "" {
		return false
	}
	if _, err := net.LookupHost(host); err != nil {
		return true
	}
	return false
}

func extractDBHost(dsn string) string {
	dsn = strings.TrimSpace(dsn)
	if dsn == "" {
		return ""
	}
	parsed, err := neturl.Parse(dsn)
	if err != nil {
		return ""
	}
	return parsed.Hostname()
}
