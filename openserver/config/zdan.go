package config

import (
	"fmt"
	"os"
)

type ZdanConfig struct {
	CloudDmappId  string `yaml:"cloudDmappId"`
	CloudDmappKey string `yaml:"cloudDmappKey"`
	CloudUserId   string `yaml:"cloudUserId"`
	UserDmappId   string `yaml:"userDmappId"`
	UserDmappKey  string `yaml:"userDmappKey"`
	ZdanHost      string `yaml:"zdanHost"`
	ZdanPort      string `yaml:"zdanPort"`
}

func (c *ZdanConfig) Check() error {

	if len(c.CloudDmappId) == 0 {
		return fmt.Errorf("invalid cloud dmapp id")
	}

	if len(c.CloudDmappKey) == 0 {
		return fmt.Errorf("invalid cloud dmapp key")
	}

	if len(c.CloudUserId) == 0 {
		return fmt.Errorf("invalid cloud user id")
	}

	if len(c.UserDmappId) == 0 {
		return fmt.Errorf("invalid user dmapp id")
	}

	if len(c.UserDmappKey) == 0 {
		return fmt.Errorf("invalid user dmapp key")
	}

	// 通过算力调度部署，使用环境变量覆盖

	host := os.Getenv("ZDAN_API_ADDRESS")
	port := os.Getenv("ZDAN_API_PORT")

	if len(host) > 0 {
		c.ZdanHost = host
	}

	if len(port) > 0 && port != "443" && port != "80" {
		c.ZdanPort = port
	}

	return nil
}

func (c ZdanConfig) Address() string {
	if len(c.ZdanPort) > 0 {
		return fmt.Sprintf("%s:%s", c.ZdanHost, c.ZdanPort)
	}

	return c.ZdanHost

}
