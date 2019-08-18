package docker

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/cogolabs/rudder/internal/config"
)

var (
	tickerInterval = 5 * time.Second
)

// WaitForImages waits for the given docker iamges and tag to be built
func WaitForImages(cfg *config.Config, tag string) error {
	ticker := time.NewTicker(tickerInterval)
	for _, cntr := range cfg.Containers {
		fmt.Printf("Waiting for %s:%s to build on %s...\n", cntr.Image, tag, cntr.Registry)
		timer := time.NewTimer(cntr.Timeout)
		for {
			ready, err := checkImage(&cntr, tag)
			if err != nil {
				return err
			}
			if ready {
				break
			}

			select {
			case <-timer.C:
				return ErrTimeout
			case <-ticker.C:
				continue
			}
		}
	}

	return nil
}

func checkImage(cntr *config.Container, tag string) (bool, error) {
	url := fmt.Sprintf("%s/v1/repositories/%s/tags/%s", cntr.Registry, cntr.Image, tag)
	res, err := http.Get(url)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()
	_, err = io.Copy(ioutil.Discard, res.Body)
	if err != nil {
		return false, err
	}

	return checkResponse(res)
}

func checkResponse(res *http.Response) (bool, error) {
	if res.StatusCode == http.StatusOK {
		return true, nil
	}
	if res.StatusCode == http.StatusNotFound {
		return false, nil
	}

	return false, &ErrBadDockerResponse{res.StatusCode, res.Request.URL.String()}
}
