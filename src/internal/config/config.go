package config

import (
	"strconv"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DBUser     string `env:"DATABASE_USER" env-default:""`
	DBHost     string `env:"DATABASE_HOST" env-default:""`
	DBName     string `env:"DATABASE_NAME" env-default:""`
	DBPassword string `env:"DATABASE_PASSWORD" env-default:""`
	DBPort     int    `env:"DATABASE_PORT" env-default:"0"`

	ServerPort int `env:"SERVER_PORT" env-required:"true"`

	Storage string `env:"STORAGE" env-required:"true"`
}

func Load() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (cfg *Config) GetDBDSN() string {
	return "host=" + cfg.DBHost +
		" port=" + strconv.Itoa(cfg.DBPort) +
		" user=" + cfg.DBUser +
		" password=" + cfg.DBPassword +
		" dbname=" + cfg.DBName +
		" sslmode=disable"
}
