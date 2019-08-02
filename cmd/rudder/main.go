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
	branch = flag.String("branch", "", "Current branch")
	tag    = flag.String("tag", "", "Current tag")
)

func main() {
	flag.Parse()
	cfg, err := config.Load()
	die(err)
	err = docker.WaitForImage(cfg, *tag)
	die(err)
}

func die(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s%s%s\n", errColor, err.Error(), endColor)
		os.Exit(1)
	}
}
