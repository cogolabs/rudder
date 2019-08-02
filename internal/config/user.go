package config

import "os"

const (
	defaultUserName = "default"
	tokenVar        = "KUBE_TOKEN"
)

// User holds information about the kubernetes user
type User struct {
	Name              string `yaml:"name"`
	Token             string `yaml:"-"`
	ClientCertificate string `yaml:"client_certificate"`
	ClientKey         string `yaml:"client_key"`
}

func (u *User) process() *User {
	if u.Name == "" {
		u.Name = defaultUserName
	}
	u.Token = os.Getenv(tokenVar)
	return u
}
