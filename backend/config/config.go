package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port string
	Mode string
}

type DatabaseConfig struct {
	Driver string
	DSN    string
}

type JWTConfig struct {
	Secret     string
	ExpireTime int64
}

var AppConfig *Config

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Warning: Config file not found: %v\n", err)
		setDefaultConfig()
	}

	AppConfig = &Config{
		Server: ServerConfig{
			Port: viper.GetString("server.port"),
			Mode: viper.GetString("server.mode"),
		},
		Database: DatabaseConfig{
			Driver: viper.GetString("database.driver"),
			DSN:    viper.GetString("database.dsn"),
		},
		JWT: JWTConfig{
			Secret:     viper.GetString("jwt.secret"),
			ExpireTime: viper.GetInt64("jwt.expire_time"),
		},
	}
}

func setDefaultConfig() {
	workDir, _ := os.Getwd()
	dbPath := fmt.Sprintf("%s/transportation.db", workDir)

	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("database.driver", "sqlite")
	viper.SetDefault("database.dsn", dbPath)
	viper.SetDefault("jwt.secret", "smart_transportation_jwt_secret_key_2024")
	viper.SetDefault("jwt.expire_time", 86400)
}
