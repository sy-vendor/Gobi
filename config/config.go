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
		Secret string
	}
	Database struct {
		Type string
		DSN  string
	}
}

var AppConfig Config

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
	AppConfig.Database.Type = viper.GetString("database.type")
	AppConfig.Database.DSN = viper.GetString("database.dsn")

	fmt.Printf("Loaded config for env: %s, port: %s, db type: %s\n", env, AppConfig.Server.Port, AppConfig.Database.Type)
}
