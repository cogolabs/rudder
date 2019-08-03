package main

import (
	"fmt"
	"os"
	"runtime/debug"
)

var (
	version = "master"
	commit  = "?"
	date    = ""
)

//nolint:gochecknoinits
func init() {
	if info, available := debug.ReadBuildInfo(); available {
		if date == "" && info.Main.Version != "(devel)" {
			version = info.Main.Version
			commit = fmt.Sprintf("(unknown, mod sum: %q)", info.Main.Sum)
			date = "(unknown)"
		}
	}
}

func printVersion() {
	fmt.Printf("rudder has version %s built from %s on %s\n", version, commit, date)
	os.Exit(0)
}
