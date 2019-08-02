package docker

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/ryantking/rudder/internal/config"
)

var (
	tickerInterval = 5 * time.Second
)

// WaitForImage waits for the given docker iamge and tag to be built
func WaitForImage(cfg *config.Config, tag string) error {
	timer := time.NewTimer(cfg.DockerTimeout)
	ticker := time.NewTicker(tickerInterval)
	defer timer.Stop()
	for {
		ready, err := checkImage(cfg, tag)
		if err != nil {
			return err
		}
		if ready {
			return nil
		}

		select {
		case <-timer.C:
			return ErrTimeout
		case <-ticker.C:
			continue
		}
	}
}

func checkImage(cfg *config.Config, tag string) (bool, error) {
	url := fmt.Sprintf("%s/v1/repositories/%s/tags/%s", cfg.DockerRegistry, cfg.DockerImage, tag)
	res, err := http.Get(url)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()
	_, err = io.Copy(ioutil.Discard, res.Body)
	if err != nil {
		return false, err
	}

	if res.StatusCode == http.StatusOK {
		return true, nil
	}
	if res.StatusCode == http.StatusNotFound {
		return false, nil
	}

	return false, &ErrBadDockerResponse{res.StatusCode, url}
}
