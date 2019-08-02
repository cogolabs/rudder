package docker

import (
	"errors"
	"fmt"
)

var (
	// ErrTimeout is thrown when
	ErrTimeout = errors.New("timed out waiting for image to build")
)

// ErrBadDockerResponse is thrown when a bad response is returned
// from the docker registry
type ErrBadDockerResponse struct {
	statusCode int
	url        string
}

func (err *ErrBadDockerResponse) Error() string {
	return fmt.Sprintf("recieved code %d from '%s'", err.statusCode, err.url)
}
