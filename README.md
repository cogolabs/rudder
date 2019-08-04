![Rudder](https://raw.githubusercontent.com/ryantking/rudder/master/static/img/logo.png)

A simple way to take control over your Kubernetes deployments.

[![Build Status](https://travis-ci.org/ryantking/rudder.svg?branch=master)](https://travis-ci.org/ryantking/rudder)
[![Test Coverage](https://api.codeclimate.com/v1/badges/e3ea6eff6537ba18ce2a/test_coverage)](https://codeclimate.com/github/ryantking/rudder/test_coverage)
[![Maintainability](https://api.codeclimate.com/v1/badges/e3ea6eff6537ba18ce2a/maintainability)](https://codeclimate.com/github/ryantking/rudder/maintainability)
[![Go Report Card](https://goreportcard.com/badge/github.com/ryantking/rudder)](https://goreportcard.com/report/github.com/ryantking/rudder)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

## Overview

Rudder is a portable tool for continuous delivery of Kubernetes applications.
Rudder can handle multiple different deployments with different conditions for
when to update. It is intended to be plugged into to any existing CI/CD
solution, such as Travis-CI that does not have out-of-the-box Kubernetes
deployment support.

Rudder is simply configured by a YAML file that gives precise control over a
project's deployments.

## Install

TODO

### Travis CI

TODO

## Quick Start

A configuration file called `.rudder.yml` is going to be searched for in the
working directory from where the script is run.

For a single cluster deployment that is deployed on tags to the master branch,
a basic configuration looks as follows.

```yaml
containers:
  - image: "library/nginx"
deployments:
  - name: prod
    branch: master
    only_tags: true
    kube_servers:
      - server: "https://mycluster.net"
    kube_deployments:
      - "deployment/nginx-deployment"
    yaml_folder: k8s
```

## Configuration

Way more can be done than just a single cluster, the following is a
configuration that adds a canary deployment to a different namespace, and
redundant deployments on two different clusters each:

```yaml
containers:
  - image: "library/nginx"
deployments:
  - name: prod
    branch: master
    only_tags: true
    tags:
      - "v*"
      - "multi-v*"
    kube_servers:
      - server: "https://k1.mycluster.net"
      - server: "https://k2.mycluster.net"
    kube_namespace: myproj
    kube_deployments:
      - "deployment/nginx-deployment"
    yaml_folder: k8s/prod
  - name: canary
    branch: master
    only_tags: true
    tags:
      - "canary-v*"
      - "multi-v*"
    kube_servers:
      - server: "https://k1.mycluster.net"
      - server: "https://k2.mycluster.net"
    kube_namespace: myproj-canary
    kube_deployments:
      - "deployment/nginx-deployment"
    yaml_folder: k8s/canary
```

### Branch and Tag

Currently, Rudder only knows how to find branch and tag information from
Travis-CI, but more CI/CD solutions will be added in the future.

Generically, the branch and tag can be added via command-line flags:

```
> rudder -help
Usage of rudder:
  -branch string
        Current branch
  -image-tag string
        Tag to use for Docker image
  -kube-config string
        Location of kube config (default "$HOME/.kube/config")
  -kubectl-version string
        Version of kubectl to use (default latest)
  -tag string
        Current tag
  -use-master
        Don't substitute 'master' with 'latest'

```

For example, in Travis, you would pass the branch and tag like such:
```yaml
after_success:
- rudder -branch=$TRAVIS_BRANCH -tag=$TRAVIS_TAG
```

The tag for the image on docker is intelligentally generated using the
following order:

1. The passed `-image-tag` flag
2. The git tag either found or provided via the `-tag` flag
3. The git branch, note that unless `-use-master` is passed, `master` will be substituted with `latest`

### Kubectl

As of now, the interaction with Kubernetes is done via wrapped calls to
`kubectl`. Installation and configuration of the config file is done
automatically. By default, the latest stable release of `kubectl` is
installed, but a version can be passed using the `-kubectl-version`
flag. `$HOME/kube/config` is also used by default, but a custom one
can be specified with the `-kube-config` flag.

### Config File

Rudder looks for a config file in the current working directory named `.rudder.yml`

Here is an example configuration file with default values filled in and explanation.

```yaml
# Configuration for the containers required for deployment
containers:
  - # Registry that hosts the image
    registry: https://index.docker.io # Docker Hub

    # Image to wait to build
    image: # REQUIRED

    # Timeout for waiting for the image to become available
    timeout: "5m"

# User configuration for interacting with Kubernetes
user:
    # Name of the user
    name: default

    # Path to the client certificate
    client_certificate:

    # Path to the client key
    client_key:

# Deployment configurations
deployments:
  - # Name of the deployment
    name: # REQUIRED

    # Branch to execute the deployment on
    branch: master

    # Only deploy on tagged releases
    only_tags: false

    # Tag patterns to match
    tags:

    # Folder of YAML Kubernetes resources to apply
    yaml_folder: "k8s"

    # Kubernetes Servers to apply the deployments to
    kube_servers:
      - # URL of the server
        server: # REQUIRED

        # Path to the certificate authority
        ca:

    # Kubernetes namespace to apply the deployments in
    kube_namespace: default

    # Kubernetes deployments to watch for rollout
    kube_deployments:
```

### Authentication

Currently, Rudder supports both token based authentication and certificate based
authenciation.

The Kubernetes token to use comes from the environment variable `$KUBE_TOKEN` and
the paths to the CA, certificate, and key are given in the config.

## Future Plans

1. Support authentication for Docker registries
2. More thorough testing of the functionality that directly communicates with kubernetes
3. Replace `kubectl` calls with the Kubernetes golang client


## Contributing

Feel free to open issues and pull requests for anything that would make this
project more useful or easier to maintain.

Would especially appreiate help from folks who use other container registers
such as GitLab Container Registry and AWS ECR.

# License Scan

[![FOSSA Status](https://app.fossa.com/api/projects/custom%2B12297%2Fgithub.com%2Fryantking%2Frudder.svg?type=large)](https://app.fossa.com/projects/custom%2B12297%2Fgithub.com%2Fryantking%2Frudder?ref=badge_large)
