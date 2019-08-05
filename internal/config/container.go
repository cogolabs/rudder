package config

import (
	"fmt"
	"time"
)

const (
	defaultTimeout  = 5 * time.Minute
	defaultRegistry = "https://index.docker.io"
)

// Container holds container configuration image
type Container struct {
	Registry   string        `json:"registry" yaml:"registry" toml:"registry"`
	Image      string        `json:"image" yaml:"image" toml:"image"`
	TimeoutStr string        `json:"timeout" yaml:"timeout" toml:"timeout"`
	Timeout    time.Duration `json:"-" yaml:"-" toml:"-"`
}

func (cntr *Container) process(n int) (*Container, error) {
	if cntr.Registry == "" {
		cntr.Registry = defaultRegistry
	}
	if cntr.Image == "" {
		return nil, &ErrMissingField{fmt.Sprintf("containers[%d].image", n)}
	}
	if cntr.TimeoutStr == "" {
		cntr.Timeout = defaultTimeout
	} else {
		timeout, err := time.ParseDuration(cntr.TimeoutStr)
		if err != nil {
			return nil, err
		}
		cntr.TimeoutStr = ""
		cntr.Timeout = timeout
	}

	return cntr, nil
}
