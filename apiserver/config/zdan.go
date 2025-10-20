package config

import (
	"fmt"
	"os"
)

type ZdanConfig struct {
	ApiServerKey string `yaml:"apiServerKey"`
	ApiServiceId string `yaml:"apiServiceId"`
	ZdanHost     string `yaml:"zdanHost"`
	ZdanPort     string `yaml:"zdanPort"`
}

func (c *ZdanConfig) Check() error {

	if len(c.ApiServerKey) == 0 {
		return fmt.Errorf("invalid api server key")
	}

	if len(c.ApiServiceId) == 0 {
		return fmt.Errorf("invalid api service id")
	}

	svcid := os.Getenv("ZDAN_RESOURCE_ID")
	host := os.Getenv("ZDAN_API_ADDRESS")
	port := os.Getenv("ZDAN_API_PORT")

	if len(host) > 0 {
		c.ZdanHost = host
	}

	if len(port) > 0 {
		c.ZdanPort = port
	}

	if len(svcid) > 0 {
		c.ApiServiceId = svcid
	}

	return nil
}

func (c ZdanConfig) Address() string {
	if len(c.ZdanPort) > 0 {
		return fmt.Sprintf("%s:%s", c.ZdanHost, c.ZdanPort)
	}

	return c.ZdanHost

}
