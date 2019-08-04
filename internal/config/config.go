package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

const configName = ".rudder.yml"

// Config represents a kubernetes configuration
type Config struct {
	User        User         `yaml:"user"`
	Containers  []Container  `yaml:"containers"`
	Deployments []Deployment `yaml:"deployments"`
}

// Load loads in all configurations from a file
func Load() (*Config, error) {
	cfg, err := readYAML()
	if err != nil {
		return nil, err
	}

	user := cfg.User.process()
	cfg.User = *user

	for i, cntr := range cfg.Containers {
		processed, err := cntr.process(i)
		if err != nil {
			return nil, err
		}
		cfg.Containers[i] = *processed
	}

	for i, dply := range cfg.Deployments {
		processed, err := dply.process(i)
		if err != nil {
			return nil, err
		}
		cfg.Deployments[i] = *processed
	}

	return cfg, nil
}

func readYAML() (*Config, error) {
	f, err := os.Open(configName)
	if err != nil {
		return nil, err
	}
	cfg := new(Config)
	err = yaml.NewDecoder(f).Decode(&cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
