package configuration

import (
	"os"
)

type Config struct {
	PostgresDB       string
	PostgresUser     string
	PostgresPassword string
	PostgresHost     string
}

func New() *Config {
	return &Config{
		PostgresDB:       os.Getenv("POSTGRES_DB"),
		PostgresUser:     os.Getenv("POSTGRES_USER"),
		PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
		PostgresHost:     os.Getenv("POSTGRES_HOST"),
	}
}
