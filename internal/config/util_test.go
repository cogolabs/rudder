package config

import "time"

var (
	testConfig = Config{
		DockerRegistry: "https://registry.server.net",
		DockerImage:    "org/repo",
		DockerTimeout:  2 * time.Minute,
		Deployments: []Deployment{
			{
				Name:            "prod",
				Branch:          "master",
				OnlyTags:        true,
				Tags:            []string{"multi-v*", "v*"},
				tagsRegex:       "^(multi-v.*)|(v.*)$",
				YAMLFolder:      "k8s",
				KubeServers:     []string{"mykubes1.test.net", "mykubes2.test.net"},
				KubeNamespace:   "myproject",
				KubeDeployments: []string{"deployment/myapi", "statefulset/myworker"},
			},
			{
				Name:            "canary",
				Branch:          "master",
				OnlyTags:        true,
				Tags:            []string{"multi-v*", "canary-v*"},
				tagsRegex:       "^(multi-v.*)|(canary-v.*)$",
				YAMLFolder:      "k8s-canary",
				KubeServers:     []string{"mykubes1.test.net", "mykubes2.test.net"},
				KubeNamespace:   "myproject",
				KubeDeployments: []string{"deployment/myapi-canary", "statefulset/myworker-canary"},
			},
			{
				Name:            "staging",
				Branch:          "master",
				YAMLFolder:      "k8s-staging",
				KubeServers:     []string{"mykubes3.test.net"},
				KubeNamespace:   "myproject",
				KubeDeployments: []string{"deployment/myapi", "statefulset/myworker"},
			},
		},
	}
)
