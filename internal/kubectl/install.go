package kubectl

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
)

const (
	kubectlBase = "https://storage.googleapis.com"
	pathBase    = "/kubernetes-release/release/%s/bin/%s/%s/kubectl"
)

var kubectlPath = "./kubectl"

// Install installs the desired version of kubectl
func Install(version string) error {
	binary, err := getBinary(version)
	if err != nil {
		return err
	}
	defer binary.Close()
	f, err := os.Create(kubectlPath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, binary)
	if err != nil {
		return err
	}

	return os.Chmod(kubectlPath, os.ModePerm)
}

func getBinary(version string) (io.ReadCloser, error) {
	path := fmt.Sprintf(pathBase, version, runtime.GOOS, runtime.GOARCH)
	url := fmt.Sprintf("%s%s", kubectlBase, path)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("could not install kubectl, received code %d", res.StatusCode)
	}

	return res.Body, nil
}

// Uninstall deletes the installed binary
func Uninstall() error {
	return os.RemoveAll(kubectlPath)
}
