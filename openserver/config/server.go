package config

import "fmt"

// ServerConfig contains server-related configuration
type ServerConfig struct {
	Port    int    `yaml:"port"`
	Host    string `yaml:"host"`
	Timeout int    `yaml:"timeout"`
}

func (c ServerConfig) Check() error {
	if c.Port <= 0 || c.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Port)
	}
	return nil
}

func (c ServerConfig) ListenAddress() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
