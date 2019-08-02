package config

import (
	"fmt"
	"strings"
)

const (
	defaultBranch     = "master"
	defaultYAMLFolder = "k8s"
	defaultNamespace  = "default"
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

func (dply *Deployment) process(n int) error {
	if dply.Name == "" {
		return &ErrMissingField{fmt.Sprintf("deployments[%d].name", n)}
	}
	if dply.Branch == "" {
		dply.Branch = defaultBranch
	}
	dply.tagRegex()
	if dply.YAMLFolder == "" {
		dply.YAMLFolder = defaultYAMLFolder
	}
	if len(dply.KubeServers) == 0 {
		return &ErrMissingField{fmt.Sprintf("deployments[%d].kube_servers", n)}
	}
	if dply.KubeNamespace == "" {
		dply.KubeNamespace = defaultNamespace
	}

	return nil
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
