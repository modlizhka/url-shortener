package config

import (
	"sync"
	"url-shortener/pkg/logging"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	Listen   Listen   `env:"LISTEN"`
	DataBase DataBase `env:"DATABASE"`
	CacheTTL int64    `env:"CACHE_TTL"`
}

type Listen struct {
	Type   string `env:"LISTEN_TYPE" env-default:"port"`
	BindIP string `env:"BIND_IP" env-default:"127.0.0.1"`
	Port   string `env:"PORT" env-default:"8080"`
}

type DataBase struct {
	Host     string `env:"DB_HOST" env-default:"postgres"`
	Port     string `env:"DB_PORT" env-default:"5432"`
	Username string `env:"DB_USERNAME" env-default:"postgres"`
	Password string `env:"DB_PASSWORD" env-default:"postgres"`
	DBName   string `env:"DB_NAME" env-default:"postgres"`
}

var instance *Config
var once sync.Once

func GetConfig(logFile, envFilePath string) *Config {
	once.Do(func() {
		logger, err := logging.GetLogger(logFile)
		if err != nil {
			panic(err)
		}
		logger.Info("read application configuration")
		instance = &Config{}
		if err := godotenv.Load(envFilePath); err != nil {
			logger.Fatal("Error loading .env file")
		}
		if err := env.Parse(instance); err != nil {
			logger.Fatalf("Failed to parse env variables: %v", err)
		}

	})
	return instance
}
