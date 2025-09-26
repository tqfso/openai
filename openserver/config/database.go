package config

import (
	"fmt"
	"os"
	"strconv"
)

type DatabaseConfig struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

func (c *DatabaseConfig) Check() error {

	if host := os.Getenv("PG_HOST"); host != "" {
		c.Host = host
	}

	if portStr := os.Getenv("PG_PORT"); portStr != "" {
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return fmt.Errorf("PG_PORT must be an integer")
		}
		c.Port = port
	}

	if user := os.Getenv("PG_USER"); user != "" {
		c.User = user
	}

	if password := os.Getenv("PG_PASSWORD"); password != "" {
		c.Password = password
	}

	if dbname := os.Getenv("PG_DBNAME"); dbname != "" {
		c.DBName = dbname
	}

	if sslmode := os.Getenv("PG_SSLMODE"); sslmode != "" {
		c.SSLMode = sslmode
	}

	return nil
}
