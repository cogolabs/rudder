package kubectl

import (
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/ryantking/rudder/internal/config"
)

// WaitForRollouts waits for a given deployment to finish all its rollouts
func WaitForRollouts(out io.Writer, dply config.Deployment) error {
	for _, kubeDeply := range dply.KubeDeployments {
		fmt.Fprintf(out, "Waiting for %s in namespace %s to rollout...\n", kubeDeply, dply.KubeNamespace)
		args := []string{kubectlPath, "rollout", "status", "-n", dply.KubeNamespace, kubeDeply}
		fmt.Fprintln(out, strings.Join(args, " "))
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdout = out
		cmd.Stderr = out
		err := cmd.Run()
		if err != nil {
			return err
		}
	}

	return nil
}
