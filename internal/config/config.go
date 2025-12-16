package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	LogLevel   string      `yaml:"log_level" env-default:"debug"`
	MSSQL      MSSQLConfig `yaml:"mssql"`
	ClickHouse CHConfig    `yaml:"clickhouse"`
}

type MSSQLConfig struct {
	DSN        string `yaml:"dsn"`
	MaxConns   int    `yaml:"max_conns"`
	RateLimit  int    `yaml:"rate_limit"`
	CDCEnabled bool   `yaml:"cdc_enabled"`
}

type CHConfig struct {
	DSN      string `yaml:"dsn"`
	MaxConns int    `yaml:"max_conns"`
}

func MustLoad(configPath string) *Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exits: %s", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
