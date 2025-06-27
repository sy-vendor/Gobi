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
		Type string
		DSN  string
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

	validateConfig(&AppConfig)

	fmt.Printf("Loaded config for env: %s, port: %s, db type: %s\n", env, AppConfig.Server.Port, AppConfig.Database.Type)
}
