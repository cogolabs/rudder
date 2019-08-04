package testutil

import (
	"io"
	"os"
)

const configName = ".rudder.yml"

// WriteConfig writes a config file to the local folder
func WriteConfig(fname string) error {
	src, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer src.Close()
	dst, err := os.Create(configName)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}

// WriteConfigTo writes a config file to a specific location
func WriteConfigTo(srcFname, dstFname string) error {
	src, err := os.Open(srcFname)
	if err != nil {
		return err
	}
	defer src.Close()
	dst, err := os.Create(dstFname)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}
