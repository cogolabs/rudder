package config

import "os"

const (
	defaultUserName = "default"
	tokenVar        = "KUBE_TOKEN"
)

// User holds information about the kubernetes user
type User struct {
	Name              string `json:"name" yaml:"name"`
	Token             string `json:"-" yaml:"-"`
	ClientCertificate string `json:"client_certificate" yaml:"client_certificate"`
	ClientKey         string `json:"client_key" yaml:"client_key"`
}

func (u *User) process() *User {
	if u.Name == "" {
		u.Name = defaultUserName
	}
	u.Token = os.Getenv(tokenVar)
	return u
}
