package config

import (
	"errors"
	"os"
)

type AppConfig struct {
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

type Config struct {
	App AppConfig
	DB  DatabaseConfig
}

func Load() (*Config, error) {
	cfg := &Config{
		App: AppConfig{
			Port: os.Getenv("APP_PORT"),
		},
		DB: DatabaseConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			Name:     os.Getenv("DB_NAME"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
		},
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	if c.App.Port == "" {
		return errors.New("APP_PORT is required")
	}

	if c.DB.Host == "" {
		return errors.New("DB_HOST is required")
	}

	if c.DB.Port == "" {
		return errors.New("DB_PORT is required")
	}

	if c.DB.Name == "" {
		return errors.New("DB_NAME is required")
	}

	if c.DB.User == "" {
		return errors.New("DB_USER is required")
	}

	if c.DB.Password == "" {
		return errors.New("DB_PASSWORD is required")
	}

	return nil
}
