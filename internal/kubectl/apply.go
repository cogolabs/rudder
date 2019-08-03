package kubectl

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const imageTagPlaceholder = "{{IMAGE_TAG}}"

// ApplyDir applies all yaml files in a directory
func ApplyDir(out io.Writer, dir, tag, kubeConfig string) error {
	return filepath.Walk(dir, processYAML(out, tag, kubeConfig))
}

func processYAML(out io.Writer, tag, kubeConfig string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) != ".yml" && filepath.Ext(path) != ".yaml" {
			return nil
		}
		err = subTag(path, tag)
		if err != nil {
			return err
		}

		args := []string{kubectlPath, "apply", "-f", path, fmt.Sprintf("--kubeconfig=%s", kubeConfig)}
		fmt.Fprintln(out, strings.Join(args, " "))
		stdout, err := exec.Command(args[0], args[1:]...).CombinedOutput()
		if err != nil {
			return err
		}

		fmt.Fprintln(out, strings.TrimSuffix(string(stdout), "\n"))
		return unstashFile(path)
	}
}

func subTag(path, tag string) error {
	err := stashFile(path)
	if err != nil {
		return err
	}
	input, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	replaced := strings.Replace(string(input), imageTagPlaceholder, tag, -1)
	return ioutil.WriteFile(path, []byte(replaced), 0644)
}

func stashFile(path string) error {
	src, err := os.Open(path)
	if err != nil {
		return err
	}
	defer src.Close()
	dstPath := fmt.Sprintf("%s.bak", path)
	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}

func unstashFile(path string) error {
	oldpath := fmt.Sprintf("%s.bak", path)
	return os.Rename(oldpath, path)
}
