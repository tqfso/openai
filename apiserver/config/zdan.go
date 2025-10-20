package config

import (
	"fmt"
	"net/url"
	"os"
)

type ZdanConfig struct {
	ApiServerKey string `yaml:"apiServerKey"`
	ApiServiceId string `yaml:"apiServiceId"`
	OpenBaseURL  string `yaml:"openBaseURL"`
}

func (c *ZdanConfig) Check() error {

	apiServiceId := os.Getenv("ZDAN_RESOURCE_ID")
	if len(apiServiceId) > 0 {
		c.ApiServiceId = apiServiceId
	}

	if len(c.ApiServiceId) == 0 {
		return fmt.Errorf("invalid api service id")
	}

	if len(c.ApiServerKey) == 0 {
		return fmt.Errorf("invalid api server key")
	}

	_, err := url.Parse(c.OpenBaseURL)
	if err != nil {
		return fmt.Errorf("invalid base URL: %w", err)
	}

	return nil
}
