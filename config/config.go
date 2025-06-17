package config

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	Type string
	DSN  string
}

type JWTConfig struct {
	Secret string
}

var DefaultConfig = Config{
	Server: ServerConfig{
		Port: "8080",
	},
	Database: DatabaseConfig{
		Type: "sqlite",
		DSN:  "gobi.db",
	},
	JWT: JWTConfig{
		Secret: "your-secret-key",
	},
}
