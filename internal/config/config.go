package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

const (
	configName      = ".rudder.yml"
	defaultTimeout  = time.Minute
	defaultRegistry = "https://index.docker.io"
)

// Config represents a kubernetes configuration
type Config struct {
	DockerRegistry   string        `yaml:"docker_registry"`
	DockerImage      string        `yaml:"docker_image"`
	DockerTimeout    time.Duration `yaml:"-"`
	DockerTimeoutStr string        `yaml:"docker_timeout"`
	Deployments      []Deployment  `yaml:"deployments"`
}

// Load loads in all configurations from a file
func Load() (*Config, error) {
	cfg, err := readYAML()
	if err != nil {
		return nil, err
	}
	err = cfg.process()
	if err != nil {
		return nil, err
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

func (cfg *Config) process() error {
	if cfg.DockerRegistry == "" {
		cfg.DockerRegistry = defaultRegistry
	}
	if cfg.DockerImage == "" {
		return &ErrMissingField{"docker_image"}
	}
	if cfg.DockerTimeoutStr == "" {
		cfg.DockerTimeout = defaultTimeout
	} else {
		timeout, err := time.ParseDuration(cfg.DockerTimeoutStr)
		if err != nil {
			return err
		}
		cfg.DockerTimeoutStr = ""
		cfg.DockerTimeout = timeout
	}

	return nil
}
