package config

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/cogolabs/rudder/internal/kubes"
	"gopkg.in/yaml.v2"
)

const (
	defaultBranch     = "master"
	defaultYAMLFolder = "k8s"
	defaultNamespace  = "default"
)

// Deployment holds the configuration info for a specific deployment
type Deployment struct {
	Name            string       `json:"name" yaml:"name" toml:"name"`
	Branch          string       `json:"branch" yaml:"branch" toml:"branch"`
	OnlyTags        bool         `json:"only_tags" yaml:"only_tags" toml:"only_tags"`
	Tags            []string     `json:"tags" yaml:"tags" toml:"tags"`
	YAMLFolder      string       `json:"yaml_folder" yaml:"yaml_folder" toml:"yaml_folder"`
	KubeServers     []KubeServer `json:"kube_servers" yaml:"kube_servers" toml:"kube_servers"`
	KubeNamespace   string       `json:"kube_namespace" yaml:"kube_namespace" toml:"kube_namespace"`
	KubeDeployments []string     `json:"kube_deployments" yaml:"kube_deployments" toml:"kube_deployments"`

	tagsRegex string
}

// KubeServer holds information about a Kubernetes server
type KubeServer struct {
	Server string `yaml:"server"`
	CA     string `yaml:"ca"`
}

func (dply *Deployment) process(n int) (*Deployment, error) {
	if dply.Name == "" {
		return nil, &ErrMissingField{fmt.Sprintf("deployments[%d].name", n)}
	}
	if dply.Branch == "" {
		dply.Branch = defaultBranch
	}
	if dply.YAMLFolder == "" {
		dply.YAMLFolder = defaultYAMLFolder
	}
	if len(dply.KubeServers) == 0 {
		return nil, &ErrMissingField{fmt.Sprintf("deployments[%d].kube_servers", n)}
	}
	for i, srv := range dply.KubeServers {
		if srv.Server == "" {
			return nil, &ErrMissingField{fmt.Sprintf("deployments[%d].kube_servers[%d].server", n, i)}
		}
	}
	if dply.KubeNamespace == "" {
		dply.KubeNamespace = defaultNamespace
	}
	if len(dply.Tags) > 0 {
		dply.genTagRegex()
	}

	return dply, nil
}

func (dply *Deployment) genTagRegex() {
	for i, tag := range dply.Tags {
		if i == 0 {
			dply.tagsRegex = fmt.Sprintf("^(%s)", strings.Replace(tag, "*", ".*", -1))
		} else {
			dply.tagsRegex = fmt.Sprintf("%s|(%s)", dply.tagsRegex, strings.Replace(tag, "*", ".*", -1))
		}
	}
	dply.tagsRegex = fmt.Sprintf("%s$", dply.tagsRegex)
}

// MakeKubesConfig makes the kubes config
func (dply *Deployment) MakeKubesConfig(user *User, configPath string, server int) error {
	config := kubes.DefaultConfig
	config.Clusters[0].Cluster.Server = dply.KubeServers[server].Server
	config.Clusters[0].Cluster.CertificateAuthority = dply.KubeServers[server].CA
	config.Contexts[0].Context.Cluster = config.Clusters[0].Name
	config.Contexts[0].Context.Namespace = dply.KubeNamespace
	config.Contexts[0].Context.User = config.Users[0].Name
	config.Users[0].User.Token = user.Token
	config.Users[0].User.ClientCertificate = user.ClientCertificate
	config.Users[0].User.ClientKey = user.ClientKey

	configPath = os.ExpandEnv(configPath)
	configDir := filepath.Dir(configPath)
	err := os.MkdirAll(configDir, os.ModePerm)
	if err != nil {
		return err
	}

	f, err := os.Create(configPath)
	if err != nil {
		return err
	}

	return yaml.NewEncoder(f).Encode(&config)
}

// ShouldDeploy returns whether or not the deployments criteria are met
func (dply *Deployment) ShouldDeploy(branch, tag string) bool {
	if dply.Branch != branch {
		return false
	}
	if dply.OnlyTags && tag == "" {
		return false
	}
	if dply.tagsRegex != "" {
		r := regexp.MustCompile(dply.tagsRegex)
		if !r.MatchString(tag) {
			return false
		}
	}

	return true
}
