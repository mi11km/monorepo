package config

import (
	"os"
	"strconv"

	"github.com/mi11km/workspaces/golang/services/template/infrastructures"
)

type Config struct {
	Debug bool
	Port  string
	MySQL infrastructures.MySQLConfig
}

func New() *Config {
	return &Config{
		Debug: GetBoolEnv("DEBUG", true),
		Port:  GetEnv("PORT", "8080"),
		MySQL: infrastructures.MySQLConfig{
			User:     os.Getenv("MYSQL_USER"),
			Password: os.Getenv("MYSQL_PASSWORD"),
			Host:     os.Getenv("MYSQL_HOST"),
			Port:     os.Getenv("MYSQL_PORT"),
			Name:     os.Getenv("MYSQL_DATABASE"),
		},
	}
}

func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func GetBoolEnv(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}
