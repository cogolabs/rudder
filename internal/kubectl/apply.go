package kubectl

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ApplyDir applies all yaml files in a directory
func ApplyDir(out io.Writer, dir, kubeConfig string) error {
	return filepath.Walk(dir, processYAML(out, kubeConfig))
}

func processYAML(out io.Writer, kubeConfig string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) != ".yml" && filepath.Ext(path) != ".yaml" {
			return nil
		}

		args := []string{kubectlPath, "apply", "-f", path, fmt.Sprintf("--kubeconfig=%s", kubeConfig)}
		fmt.Fprintln(out, strings.Join(args, " "))
		stdout, err := exec.Command(args[0], args[1:]...).CombinedOutput()
		if err != nil {
			return err
		}

		fmt.Fprintln(out, strings.TrimSuffix(string(stdout), "\n"))
		return nil
	}
}
