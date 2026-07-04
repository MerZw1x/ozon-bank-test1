package config

import (
	"strconv"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DBUser     string `env:"DATABASE_USER" env-required:"true"`
	DBHost     string `env:"DATABASE_HOST" env-required:"true"`
	DBName     string `env:"DATABASE_NAME" env-required:"true"`
	DBPassword string `env:"DATABASE_PASSWORD" env-required:"true"`
	DBPort     int    `env:"DATABASE_PORT" env-required:"true"`

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
