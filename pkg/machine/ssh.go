package machine

import (
	"fmt"
	"os/exec"
	"strconv"

	"github.com/sirupsen/logrus"
)

// CommonSSH is a common function for ssh'ing to a podman machine using system-connections
// and a port
// TODO This should probably be taught about an machineconfig to reduce input
func CommonSSH(username, identityPath, name string, sshPort int, inputArgs []string) error {
	return commonSSH(username, identityPath, name, sshPort, inputArgs, false)
}

func CommonSSHSilent(username, identityPath, name string, sshPort int, inputArgs []string) error {
	return commonSSH(username, identityPath, name, sshPort, inputArgs, true)
}

func commonSSH(username, identityPath, name string, sshPort int, inputArgs []string, silent bool) error {
	sshDestination := username + "@localhost"
	port := strconv.Itoa(sshPort)
	interactive := true

	args := []string{"-i", identityPath, "-p", port, sshDestination,
		"-o", "IdentitiesOnly=yes",
		"-o", "StrictHostKeyChecking=no", "-o", "LogLevel=ERROR", "-o", "SetEnv=LC_ALL="}
	if len(inputArgs) > 0 {
		interactive = false
		args = append(args, inputArgs...)
	} else {
		// ensure we have a tty
		args = append(args, "-t")
		fmt.Printf("Connecting to vm %s. To close connection, use `~.` or `exit`\n", name)
	}

	cmd := exec.Command("ssh", args...)
	logrus.Debugf("Executing: ssh %v\n", args)

	if !silent {
		if err := setupIOPassthrough(cmd, interactive); err != nil {
			return err
		}
	}

	return cmd.Run()
}
