package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const (
	configBase = ".rudder"

	extJSON = ".json"
	extYML  = ".yml"
	extYAML = ".yaml"
)

var (
	configFormats = []string{"yaml", "json"}

	// ErrConfigNotFound is thrown when a configuration file is not found
	ErrConfigNotFound = errors.New("no configuration file could be found")
)

// Config represents a kubernetes configuration
type Config struct {
	User        User         `json:"user"`
	Containers  []Container  `json:"containers"`
	Deployments []Deployment `json:"deployments"`
}

type decoder interface {
	Decode(interface{}) error
}

// Load loads in all configurations from a file
func Load() (*Config, error) {
	cfg, err := readConfig()
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

func readConfig() (*Config, error) {
	f, path, err := openConfig()
	if err != nil {
		return nil, err
	}

	var d decoder
	switch filepath.Ext(path) {
	case extJSON:
		d = json.NewDecoder(f)
	case extYAML, extYML:
		d = yaml.NewDecoder(f)
	default:
		return nil, fmt.Errorf("unsupported config format: %s", filepath.Ext(path))
	}
	cfg := new(Config)
	err = d.Decode(&cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func openConfig() (io.ReadCloser, string, error) {
	matches, err := filepath.Glob(fmt.Sprintf("./%s.*", configBase))
	if err != nil {
		return nil, "", err
	}
	if len(matches) == 0 {
		return nil, "", ErrConfigNotFound
	}
	path := matches[0]
	f, err := os.Open(path)
	if err != nil {
		return nil, "", err
	}

	return f, path, nil
}
