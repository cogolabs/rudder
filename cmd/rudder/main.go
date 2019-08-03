package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/ryantking/rudder/internal/config"
	"github.com/ryantking/rudder/internal/docker"
	"github.com/ryantking/rudder/internal/kubectl"
)

const (
	errColor = "\033[91m"
	endColor = "\033[0m"

	travisBranchVar = "TRAVIS_BRANCH"
	travisTagVar    = "TRAVIS_TAG"

	kubectlVersionURL = "https://storage.googleapis.com/kubernetes-release/release/stable.txt"
)

var (
	passedBranch         = flag.String("branch", "", "Current branch")
	passedTag            = flag.String("tag", "", "Current tag")
	kubeConfig           = flag.String("kube-config", "$HOME/.kube/config", "Location of kube config")
	passedKubectlVersion = flag.String("kubectl-version", "", "Version of kubectl to use (default latest)")
)

func main() {
	flag.Parse()
	cfg, err := config.Load()
	die(err)
	branch, tag := branchAndTag()
	fmt.Printf("Waiting for %s:%s to build on %s...\n", cfg.DockerImage, tag, cfg.DockerRegistry)
	err = docker.WaitForImage(cfg, tag)
	die(err)
	kctlV, err := kubectlVersion()
	die(err)
	err = kubectl.Install(kctlV)
	die(err)

	for _, dply := range cfg.Deployments {
		if !dply.ShouldDeploy(branch, tag) {
			fmt.Printf("%s does not have its requirements met, skipping...\n", dply.Name)
			continue
		}

		for i := range dply.KubeServers {
			err = dply.MakeKubesConfig(&cfg.User, *kubeConfig, i)
			die(err)
			fmt.Printf("Deploying %s to %s on %s\n", dply.Name, dply.KubeNamespace, dply.KubeServers[i].Server)
			err = kubectl.ApplyDir(os.Stdout, dply.YAMLFolder, *kubeConfig)
		}
	}

	err = kubectl.Uninstall()
	die(err)
}

func die(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s%s%s\n", errColor, err.Error(), endColor)
		os.Exit(1)
	}
}

func branchAndTag() (string, string) {
	branch := os.Getenv(travisBranchVar)
	if branch != "" {
		return branch, os.Getenv(travisTagVar)
	}

	return *passedBranch, *passedTag
}

func kubectlVersion() (string, error) {
	if *passedKubectlVersion != "" {
		return *passedKubectlVersion, nil
	}

	res, err := http.Get(kubectlVersionURL)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("could not find kubectl latest version: code %d", res.StatusCode)
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return strings.TrimSuffix(string(b), "\n"), nil
}
