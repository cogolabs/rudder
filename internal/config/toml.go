package config

import (
	"io"

	"github.com/BurntSushi/toml"
)

type tomlDecoder struct {
	r io.Reader
}

func (td *tomlDecoder) Decode(v interface{}) error {
	_, err := toml.DecodeReader(td.r, v)
	return err
}
