package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"log/slog"
)

type Env struct {
	DB_NAME     string `env:"DB_NAME"`
	DB_USERNAME string `env:"DB_USERNAME"`
	DB_PASSWORD string `env:"DB_PASSWORD"`
	DB_PORT     int    `env:"DB_PORT"`
	DB_HOST     string `env:"DB_HOST"`
	IpAddress   string `env:"IP_ADDRESS"`
	API_PORT    int    `env:"API_PORT"`

	Environment string `env:"ENVIRONMENT"`
}

type Config struct {
	Env    Env `env:"ENVIRONMENT"`
	Client *pgxpool.Pool
}

var config Config

func GetConfig() *Config {
	config.Env = *GetEnv()

	return &config
}

func GetEnv() *Env {
	err := godotenv.Load()
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
