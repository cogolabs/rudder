package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ryantking/rudder/internal/config"
)

var (
	branch = flag.String("branch", "", "Current branch")
	tag    = flag.String("tag", "", "Current tag")
)

func main() {
	flag.Parse()
	_, err := config.Load()
	die(err)
}

func die(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
}
