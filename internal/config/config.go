package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
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

type AuthConfig struct {
	JWTSecret string
}

type BootstrapConfig struct {
	AdminEmail    string
	AdminPassword string
}

type Config struct {
	App       AppConfig
	DB        DatabaseConfig
	Auth      AuthConfig
	Bootstrap BootstrapConfig
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

		Auth: AuthConfig{
			JWTSecret: os.Getenv("JWT_SECRET"),
		},

		Bootstrap: BootstrapConfig{
			AdminEmail:    os.Getenv("ADMIN_EMAIL"),
			AdminPassword: os.Getenv("ADMIN_PASSWORD"),
		},
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func LoadEnv(path string) error {
	if err := godotenv.Load(path); err != nil {
		return fmt.Errorf("load env %s: %w", path, err)
	}

	return nil
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

	if c.Auth.JWTSecret == "" {
		return errors.New("JWT_SECRET is required")
	}

	return nil
}
