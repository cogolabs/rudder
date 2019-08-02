package config

import "time"

var (
	testConfig = Config{
		DockerImage:   "registry.test.net/org/repo",
		DockerTimeout: 2 * time.Minute,
		Deployments: []*Deployment{
			{
				Name:            "prod",
				Branch:          "master",
				OnlyTags:        true,
				Tags:            []string{"multi-v*", "v*"},
				TagsRegex:       "^(multi-v.*)|(v.*)$",
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
				TagsRegex:       "^(multi-v.*)|(canary-v.*)$",
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
