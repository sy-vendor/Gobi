package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port string
	}
	JWT struct {
		Secret          string
		ExpirationHours int
	}
	Database struct {
		Type            string
		DSN             string
		MaxOpenConns    int
		MaxIdleConns    int
		ConnMaxLifetime int
		ConnectionPool  struct {
			MaxOpenConns    int `yaml:"max_open_conns"`
			MaxIdleConns    int `yaml:"max_idle_conns"`
			ConnMaxLifetime int `yaml:"conn_max_lifetime"`
			ConnMaxIdleTime int `yaml:"conn_max_idle_time"`
		} `yaml:"connection_pool"`
	}
	Cache struct {
		Enabled  bool
		TTL      int
		MaxSize  int `yaml:"max_size"`
		Strategy struct {
			SimpleQueryTTL  int `yaml:"simple_query_ttl"`
			ComplexQueryTTL int `yaml:"complex_query_ttl"`
			MaxCacheSize    int `yaml:"max_cache_size"`
		}
	}
}

var AppConfig Config

func validateConfig(cfg *Config) {
	if cfg.Server.Port == "" {
		panic("[Config] server.port is required in config.yaml or environment variable")
	}
	if cfg.JWT.Secret == "" {
		panic("[Config] jwt.secret is required in config.yaml or environment variable")
	}
	if cfg.Database.Type == "" {
		panic("[Config] database.type is required in config.yaml or environment variable")
	}
	if cfg.Database.DSN == "" {
		panic("[Config] database.dsn is required in config.yaml or environment variable")
	}
	secret := os.Getenv("DATA_SOURCE_SECRET")
	if len(secret) != 32 {
		panic("[Config] DATA_SOURCE_SECRET must be 32 characters (256 bits) for AES-256 encryption. Current length: " + fmt.Sprint(len(secret)))
	}
}

func LoadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	_ = viper.ReadInConfig()

	env := os.Getenv("GOBI_ENV")
	if env == "" {
		env = "default"
	}

	sub := viper.Sub(env)
	if sub != nil {
		viper.MergeConfigMap(sub.AllSettings())
	}

	viper.AutomaticEnv()

	AppConfig.Server.Port = viper.GetString("server.port")
	AppConfig.JWT.Secret = viper.GetString("jwt.secret")
	AppConfig.JWT.ExpirationHours = viper.GetInt("jwt.expiration_hours")
	AppConfig.Database.Type = viper.GetString("database.type")
	AppConfig.Database.DSN = viper.GetString("database.dsn")
	AppConfig.Database.MaxOpenConns = viper.GetInt("database.max_open_conns")
	AppConfig.Database.MaxIdleConns = viper.GetInt("database.max_idle_conns")
	AppConfig.Database.ConnMaxLifetime = viper.GetInt("database.conn_max_lifetime")
	AppConfig.Database.ConnectionPool.MaxOpenConns = viper.GetInt("database.connection_pool.max_open_conns")
	AppConfig.Database.ConnectionPool.MaxIdleConns = viper.GetInt("database.connection_pool.max_idle_conns")
	AppConfig.Database.ConnectionPool.ConnMaxLifetime = viper.GetInt("database.connection_pool.conn_max_lifetime")
	AppConfig.Database.ConnectionPool.ConnMaxIdleTime = viper.GetInt("database.connection_pool.conn_max_idle_time")
	AppConfig.Cache.Enabled = viper.GetBool("cache.enabled")
	AppConfig.Cache.TTL = viper.GetInt("cache.ttl")
	AppConfig.Cache.MaxSize = viper.GetInt("cache.max_size")
	AppConfig.Cache.Strategy.SimpleQueryTTL = viper.GetInt("cache.strategy.simple_query_ttl")
	AppConfig.Cache.Strategy.ComplexQueryTTL = viper.GetInt("cache.strategy.complex_query_ttl")
	AppConfig.Cache.Strategy.MaxCacheSize = viper.GetInt("cache.strategy.max_cache_size")

	validateConfig(&AppConfig)

	fmt.Printf("Loaded config for env: %s, port: %s, db type: %s\n", env, AppConfig.Server.Port, AppConfig.Database.Type)
}
