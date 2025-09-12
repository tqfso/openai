package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var (
	config Config
)

type Config struct {
	Server ServerConfig `yaml:"server"`
}

func (c Config) Check() error {
	if err := c.Server.Check(); err != nil {
		return err
	}
	return nil
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
