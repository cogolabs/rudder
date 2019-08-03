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

	masterTag = "master"
	latestTag = "latest"

	kubectlVersionURL = "https://storage.googleapis.com/kubernetes-release/release/stable.txt"
)

var (
	passedBranch         = flag.String("branch", "", "Current branch")
	passedTag            = flag.String("tag", "", "Current tag")
	useMaster            = flag.Bool("use-master", false, "Don't substitute 'master' with 'latest'")
	passedImageTag       = flag.String("image-tag", "", "Tag to use for Docker image ")
	kubeConfig           = flag.String("kube-config", "$HOME/.kube/config", "Location of kube config")
	passedKubectlVersion = flag.String("kubectl-version", "", "Version of kubectl to use (default latest)")
	versionCmd           = flag.Bool("version", false, "Show version information and exit")
)

func main() {
	flag.Parse()
	if *versionCmd {
		printVersion()
	}

	cfg, err := config.Load()
	die(err)
	branch, tag := branchAndTag()
	imageTag := imageTag(branch, tag)

	deployments := make([]config.Deployment, 0, len(cfg.Deployments))
	for _, dply := range cfg.Deployments {
		if !dply.ShouldDeploy(branch, tag) {
			fmt.Printf("%s does not have its requirements met, skipping...\n", dply.Name)
			continue
		}
		deployments = append(deployments, dply)
	}

	if len(deployments) == 0 {
		fmt.Println("no deployments found to update, exiting...")
		return
	}

	fmt.Printf("Waiting for %s:%s to build on %s...\n", cfg.DockerImage, imageTag, cfg.DockerRegistry)
	err = docker.WaitForImage(cfg, imageTag)
	die(err)
	kctlV, err := kubectlVersion()
	die(err)
	err = kubectl.Install(kctlV)
	die(err)

	for _, dply := range deployments {
		for i := range dply.KubeServers {
			err = dply.MakeKubesConfig(&cfg.User, *kubeConfig, i)
			die(err)
			fmt.Printf("Deploying %s to %s on %s\n", dply.Name, dply.KubeNamespace, dply.KubeServers[i].Server)
			err = kubectl.ApplyDir(os.Stdout, dply.YAMLFolder, imageTag, *kubeConfig)
			die(err)
			err = kubectl.WaitForRollouts(os.Stdout, dply)
			die(err)
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

func imageTag(branch, tag string) string {
	if *passedImageTag != "" {
		return *passedImageTag
	}
	if tag != "" {
		return tag
	}
	if branch == masterTag && !*useMaster {
		return latestTag
	}

	return branch
}

func kubectlVersion() (string, error) {
	if *passedKubectlVersion != "" {
		return *passedKubectlVersion, nil
	}

	return getLatestKubectl()
}

func getLatestKubectl() (string, error) {
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
