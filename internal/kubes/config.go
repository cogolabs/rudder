package kubes

import (
	"fmt"
	"os"

	"github.com/ryantking/rudder/internal/config"
	"gopkg.in/yaml.v2"
)

const (
	tokenVar   = "KUBE_TOKEN"
	configName = "config"
)

// Config is the Kubernetes config
type Config struct {
	APIVersion     string            `yaml:"apiVersion"`
	Clusters       []Cluster         `yaml:"clusters"`
	Contexts       []Context         `yaml:"contexts"`
	CurrentContext string            `yaml:"current-context"`
	Kind           string            `yaml:"kind"`
	Preferences    map[string]string `yaml:"preferences"`
	Users          []User            `yaml:"users"`
}

// Cluster is a kubernetes cluster config
type Cluster struct {
	Name    string `yaml:"name"`
	Cluster struct {
		Server string `yaml:"server"`
	} `yaml:"cluster"`
}

// Context is a kubernetes context config
type Context struct {
	Name    string `yaml:"name"`
	Context struct {
		Cluster   string `yaml:"cluster"`
		Namespace string `yaml:"namespace"`
		User      string `yaml:"user"`
	} `yaml:"context"`
}

// User is a kubernetes user config
type User struct {
	Name string `yaml:"name"`
	User struct {
		Token string `yaml:"token"`
	} `yaml:"user"`
}

var defaultConfig = Config{
	APIVersion:     "v1",
	Clusters:       []Cluster{{Name: "cluster"}},
	Contexts:       []Context{{Name: "cluster"}},
	CurrentContext: "cluster",
	Kind:           "config",
	Preferences:    make(map[string]string),
	Users:          []User{{Name: "default"}},
}

// MakeConfig makes a configuration for a deployment
func MakeConfig(configDir string, dply config.Deployment, serverIndex int) error {
	config := defaultConfig
	config.Clusters[0].Cluster.Server = dply.KubeServers[serverIndex]
	config.Contexts[0].Context.Cluster = config.Clusters[0].Name
	config.Contexts[0].Context.Namespace = dply.KubeNamespace
	config.Contexts[0].Context.User = config.Users[0].Name
	config.Users[0].User.Token = os.Getenv(tokenVar)

	err := os.MkdirAll(configDir, os.ModePerm)
	if err != nil {
		return err
	}

	f, err := os.Create(fmt.Sprintf("%s/%s", configDir, configName))
	if err != nil {
		return err
	}

	return yaml.NewEncoder(f).Encode(&config)
}
