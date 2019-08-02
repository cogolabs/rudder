package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/ryantking/rudder/internal/kubes"
	"gopkg.in/yaml.v2"
)

const (
	defaultBranch     = "master"
	defaultYAMLFolder = "k8s"
	defaultNamespace  = "default"

	tokenVar    = "KUBE_TOKEN"
	kubesConfig = "config"
)

// Deployment holds the configuration info for a specific deployment
type Deployment struct {
	Name            string   `yaml:"name"`
	Branch          string   `yaml:"branch"`
	OnlyTags        bool     `yaml:"only_tags"`
	Tags            []string `yaml:"tags"`
	TagsRegex       string   `yaml:"-"`
	YAMLFolder      string   `yaml:"yaml_folder"`
	KubeServers     []string `yaml:"kube_servers"`
	KubeNamespace   string   `yaml:"kube_namespace"`
	KubeDeployments []string `yaml:"kube_deployments"`
}

func (dply *Deployment) process(n int) (*Deployment, error) {
	if dply.Name == "" {
		return nil, &ErrMissingField{fmt.Sprintf("deployments[%d].name", n)}
	}
	if dply.Branch == "" {
		dply.Branch = defaultBranch
	}
	dply.tagRegex()
	if dply.YAMLFolder == "" {
		dply.YAMLFolder = defaultYAMLFolder
	}
	if len(dply.KubeServers) == 0 {
		return nil, &ErrMissingField{fmt.Sprintf("deployments[%d].kube_servers", n)}
	}
	if dply.KubeNamespace == "" {
		dply.KubeNamespace = defaultNamespace
	}

	return dply, nil
}

func (dply *Deployment) tagRegex() {
	if len(dply.Tags) == 0 {
		return
	}

	for i, tag := range dply.Tags {
		if i == 0 {
			dply.TagsRegex = fmt.Sprintf("^(%s)", strings.Replace(tag, "*", ".*", -1))
		} else {
			dply.TagsRegex = fmt.Sprintf("%s|(%s)", dply.TagsRegex, strings.Replace(tag, "*", ".*", -1))
		}
	}
	dply.TagsRegex = fmt.Sprintf("%s$", dply.TagsRegex)
}

// MakeKubesConfig makes the kubes config
func (dply *Deployment) MakeKubesConfig(configDir string, server int) error {
	config := kubes.DefaultConfig
	config.Clusters[0].Cluster.Server = dply.KubeServers[server]
	config.Contexts[0].Context.Cluster = config.Clusters[0].Name
	config.Contexts[0].Context.Namespace = dply.KubeNamespace
	config.Contexts[0].Context.User = config.Users[0].Name
	config.Users[0].User.Token = os.Getenv(tokenVar)

	err := os.MkdirAll(configDir, os.ModePerm)
	if err != nil {
		return err
	}

	f, err := os.Create(fmt.Sprintf("%s/%s", configDir, kubesConfig))
	if err != nil {
		return err
	}

	return yaml.NewEncoder(f).Encode(&config)
}
