package config

import "fmt"

type ZdanConfig struct {
	CloudDmappId  string `yaml:"cloudDmappId"`
	CloudDmappKey string `yaml:"cloudDmappKey"`
	UserDmappId   string `yaml:"userDmappId"`
	UserDmappKey  string `yaml:"userDmappKey"`
}

func (c ZdanConfig) Check() error {

	if len(c.CloudDmappId) == 0 {
		return fmt.Errorf("invalid cloud dmapp id")
	}

	if len(c.CloudDmappKey) == 0 {
		return fmt.Errorf("invalid cloud dmapp key")
	}

	if len(c.UserDmappId) == 0 {
		return fmt.Errorf("invalid user dmapp id")
	}

	if len(c.UserDmappKey) == 0 {
		return fmt.Errorf("invalid user dmapp key")
	}

	return nil
}
