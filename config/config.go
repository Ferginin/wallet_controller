package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"log/slog"
)

type Env struct {
	DbName     string `env:"DB_NAME"`
	DbUsername string `env:"DB_USERNAME"`
	DbPassword string `env:"DB_PASSWORD"`
	DbPort     int    `env:"DB_PORT"`
	DbHost     string `env:"DB_HOST"`
	IpAddress  string `env:"IP_ADDRESS"`
	ApiPort    int    `env:"API_PORT"`

	Environment string `env:"ENVIRONMENT"`
}

type Config struct {
	Env    Env
	Client *pgxpool.Pool
}

var config Config

func GetConfig() *Config {
	config.Env = *GetEnv()

	return &config
}

func GetEnv() *Env {
	err := godotenv.Load("config.env")
	if err != nil {
		slog.Warn("Error loading .env file")
	}

	var cfg Env
	err = env.Parse(&cfg)
	if err != nil {
		slog.Error("Error parsing .env file:", err.Error(), nil)
		panic(err)
	}

	return &cfg
}
