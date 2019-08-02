package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ryantking/rudder/internal/config"
	"github.com/ryantking/rudder/internal/docker"
)

const (
	errColor = "\033[91m"
	endColor = "\033[0m"
)

var (
	branch     = flag.String("branch", "", "Current branch")
	tag        = flag.String("tag", "", "Current tag")
	kubeConfig = flag.String("kube-config", "$HOME/.kube/config", "Location of kube config")
)

func main() {
	flag.Parse()
	cfg, err := config.Load()
	die(err)
	fmt.Printf("Waiting for %s:%s to builld on %s...\n", cfg.DockerImage, *tag, cfg.DockerRegistry)
	err = docker.WaitForImage(cfg, *tag)
	die(err)

	for _, dply := range cfg.Deployments {
		if !dply.ShouldDeploy(*branch, *tag) {
			fmt.Printf("%s does not have its requirements met, skipping...\n", dply.Name)
			continue
		}
		for i := range dply.KubeServers {
			err = dply.MakeKubesConfig(*kubeConfig, i)
			die(err)
			fmt.Printf("Deploying %s to %s on %s\n", dply.Name, dply.KubeNamespace, dply.KubeServers[i])
		}
	}
}

func die(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s%s%s\n", errColor, err.Error(), endColor)
		os.Exit(1)
	}
}
