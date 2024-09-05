package config

import (
	"fmt"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	DBConfig     DBConfig
	ServerConfig ServerConfig
	SMTPConfig   SMTPConfig
}

type DBConfig struct {
	PgUser     string `env:"PGUSER"`
	PgPassword string `env:"PGPASSWORD"`
	PgHost     string `env:"PGHOST"`
	PgPort     uint16 `env:"PGPORT"`
	PgDatabase string `env:"PGDATABASE"`
	PgSSLMode  string `env:"PGSSLMODE"`
}

type ServerConfig struct {
	ServerMode  string `env:"ENVIRONMENT" envDefault:"debug"`
	ServerPort  string `env:"HTTP_PORT" envDefault:"8080"`
	LogFilePath string `env:"LOG_FILE_PATH"`
	SecretKey   string `env:"SECRET_KEY,notEmpty"`
}

type SMTPConfig struct {
	SmtpHost       string `env:"SMTP_HOST"`
	SmtpPort       string `env:"SMTP_PORT"`
	SenderEmail    string `env:"SENDER_EMAIL"`
	SenderPassword string `env:"SENDER_PASSWORD"`
}

func NewConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("failed to load .env file: %w", err)
	}
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config from environment variables: %w", err)
	}

	return cfg, nil
}
