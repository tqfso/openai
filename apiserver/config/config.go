package config

import (
	"fmt"
	"os"
	"path/filepath"

	"common/logger"

	"gopkg.in/yaml.v3"
)

var (
	config Config
)

type Config struct {
	Log    logger.Config `yaml:"log"`
	Server ServerConfig  `yaml:"server"`
}

type ServerConfig struct {
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
}

func (c *Config) Check() error {

	if c.Log.Level == "" {
		c.Log = logger.DefaultConfig()
	}

	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}

	return nil
}

func GetConfig() *Config {
	return &config
}

func GetLog() *logger.Config {
	return &config.Log
}

func GetServer() *ServerConfig {
	return &config.Server
}

func Load(filename string) error {

	if !filepath.IsAbs(filename) {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		filename = filepath.Join(cwd, filename)
	}

	if _, err := os.Stat(filename); err != nil {
		return err
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return err
	}

	if err := config.Check(); err != nil {
		return err
	}

	return nil
}

func Save(filename string) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return err
	}

	return nil
}
