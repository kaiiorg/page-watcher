package config

import (
	"github.com/hashicorp/hcl/v2/hclsimple"
)

type Config struct {
	Pages []*Page `hcl:"page,block"`
}

func LoadFromFile(filepath string) (*Config, error) {
	config := &Config{}

	err := hclsimple.DecodeFile(filepath, nil, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
