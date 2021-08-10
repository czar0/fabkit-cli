package docker

import (
	"fmt"

	"github.com/czar0/fabkit-cli/pkg/shell"
)

// CheckServerRunning checks the docker daemon is running in the background or returns an error
func CheckServerRunning() error {
	_, _, err := shell.PipeCommands(
		shell.NewCommand(`docker info --format '{{json .ServerErrors}}'`),
		shell.NewCommand(`grep 'null'`))
	if err != nil {
		return fmt.Errorf("error connecting to docker socket. Check your docker daemon (docker or dockerd) is running")
	}

	return nil
}
