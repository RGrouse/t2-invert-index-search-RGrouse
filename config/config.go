package config

import (
	"github.com/caarlos0/env"
)

type Config struct {
	Listen   		string `env:"LISTEN" envDefault:"127.0.0.1:8080"`
	SearchingFolder string `env:"SDIR" envDefault:"./search"`
	LogLevel 		string `env:"LOG_LEVEL" envDefault:"info"`
	PgSQL    		string `env:"PGSQL" envDefault:"postgres://rgrouse:qwedfgq4@51.15.112.136:5432/rgrouse_db?sslmode=disable"`
	ModelType		string `env:"MODEL" envDefault:"MAP"` //DB,MAP
}

func Load() (Config, error) {
	cfg := Config{}
	err := env.Parse(&cfg)
	return cfg, err
}
