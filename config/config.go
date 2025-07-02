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
			SimpleQueryTTL      int     `yaml:"simple_query_ttl"`
			ComplexQueryTTL     int     `yaml:"complex_query_ttl"`
			MaxCacheSize        int     `yaml:"max_cache_size"`
			HotCacheEnabled     bool    `yaml:"hot_cache_enabled"`
			HotCacheRatio       float64 `yaml:"hot_cache_ratio"`
			PromotionThreshold  int     `yaml:"promotion_threshold"`
			BusinessHoursStart  int     `yaml:"business_hours_start"`
			BusinessHoursEnd    int     `yaml:"business_hours_end"`
			AdaptiveTTL         bool    `yaml:"adaptive_ttl"`
			CacheWarmup         bool    `yaml:"cache_warmup"`
			MaintenanceInterval int     `yaml:"maintenance_interval"`
			EvictionPolicy      string  `yaml:"eviction_policy"`
			CompressionEnabled  bool    `yaml:"compression_enabled"`
			MetricsEnabled      bool    `yaml:"metrics_enabled"`
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
	AppConfig.Cache.Strategy.HotCacheEnabled = viper.GetBool("cache.strategy.hot_cache_enabled")
	AppConfig.Cache.Strategy.HotCacheRatio = viper.GetFloat64("cache.strategy.hot_cache_ratio")
	AppConfig.Cache.Strategy.PromotionThreshold = viper.GetInt("cache.strategy.promotion_threshold")
	AppConfig.Cache.Strategy.BusinessHoursStart = viper.GetInt("cache.strategy.business_hours_start")
	AppConfig.Cache.Strategy.BusinessHoursEnd = viper.GetInt("cache.strategy.business_hours_end")
	AppConfig.Cache.Strategy.AdaptiveTTL = viper.GetBool("cache.strategy.adaptive_ttl")
	AppConfig.Cache.Strategy.CacheWarmup = viper.GetBool("cache.strategy.cache_warmup")
	AppConfig.Cache.Strategy.MaintenanceInterval = viper.GetInt("cache.strategy.maintenance_interval")
	AppConfig.Cache.Strategy.EvictionPolicy = viper.GetString("cache.strategy.eviction_policy")
	AppConfig.Cache.Strategy.CompressionEnabled = viper.GetBool("cache.strategy.compression_enabled")
	AppConfig.Cache.Strategy.MetricsEnabled = viper.GetBool("cache.strategy.metrics_enabled")

	validateConfig(&AppConfig)

	fmt.Printf("Loaded config for env: %s, port: %s, db type: %s\n", env, AppConfig.Server.Port, AppConfig.Database.Type)
}
