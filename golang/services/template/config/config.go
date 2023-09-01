package config

import (
	"os"

	"github.com/mi11km/workspaces/golang/services/template/infrastructures"
)

type Config struct {
	Port  string
	MySQL infrastructures.MySQLConfig
}

func New() *Config {
	return &Config{
		Port: os.Getenv("PORT"),
		MySQL: infrastructures.MySQLConfig{
			User:     os.Getenv("MYSQL_USER"),
			Password: os.Getenv("MYSQL_PASSWORD"),
			Host:     os.Getenv("MYSQL_HOST"),
			Port:     os.Getenv("MYSQL_PORT"),
			Name:     os.Getenv("MYSQL_DATABASE"),
		},
	}
}
