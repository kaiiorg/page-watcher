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

	err = config.Validate()
	if err != nil {
		return nil, err
	}

	return config, nil
}

// Validate performs validation work on configuration
func (c *Config) Validate() error {
	for _, p := range c.Pages {
		err := p.ValidateNormalize()
		if err != nil {
			return err
		}
	}
	return nil
}
