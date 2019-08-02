package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

const (
	configName     = ".rudder.yml"
	defaultTimeout = time.Minute
)

// Config represents a kubernetes configuration
type Config struct {
	DockerImage      string        `yaml:"docker_image"`
	DockerTimeout    time.Duration `yaml:"-"`
	DockerTimeoutStr string        `yaml:"docker_timeout"`
	Deployments      []*Deployment `yaml:"deployments"`
}

// Load loads in all configurations from a file
func Load() (*Config, error) {
	f, err := os.Open(configName)
	if err != nil {
		return nil, err
	}
	cfg := new(Config)
	err = yaml.NewDecoder(f).Decode(&cfg)
	if err != nil {
		return nil, err
	}
	if cfg.DockerImage == "" {
		return nil, &ErrMissingField{"docker_image"}
	}
	if cfg.DockerTimeoutStr == "" {
		cfg.DockerTimeout = defaultTimeout
	} else {
		timeout, err := time.ParseDuration(cfg.DockerTimeoutStr)
		if err != nil {
			return nil, err
		}
		cfg.DockerTimeoutStr = ""
		cfg.DockerTimeout = timeout
	}
	for i, dply := range cfg.Deployments {
		err := dply.process(i)
		if err != nil {
			return nil, err
		}
	}

	return cfg, nil
}
