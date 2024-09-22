package config

import (
	"github.com/caarlos0/env/v11"
)

type (
	Config struct {
		TgToken string `env:"TG_TOKEN,notEmpty"`
		Postgres
		Redis
	}

	Postgres struct {
		Host     string `env:"PSQL_HOST" envDefault:"localhost"`
		Port     uint16 `env:"PSQL_PORT" envDefault:"5432"`
		User     string `env:"PSQL_USER"`
		Password string `env:"PSQL_PASSWORD"`
		DbName   string `env:"PSQL_DB" envDefault:"form_bot"`
		SslMode  string `env:"PSQL_SSL_MODE" envDefault:"disable"`
	}

	Redis struct {
		Host     string `env:"REDIS_HOST" envDefault:"localhost"`
		Port     uint16 `env:"REDIS_PORT" envDefault:"6379"`
		User     string `env:"REDIS_USER"`
		Password string `env:"REDIS_PASSWORD"`
		Db       int    `env:"REDIS_DB"`
	}
)

func SetupConfig() *Config {
	cfg := new(Config)
	if err := env.Parse(cfg); err != nil {
		panic(err)
	}
	return cfg
}
