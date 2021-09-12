package network

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/czar0/fabkit-cli/common/config"
	"github.com/czar0/fabkit-cli/internal/docker"
	"github.com/czar0/fabkit-cli/internal/spinner"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

// NewCmdNetworkInstall implements the method to install all the necessary dependencies for running a network
// TODO: Implement flag (-v, --version) to allow user to provide in input a specific version to install
func NewCmdNetworkInstall() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "install",
		Aliases: []string{"i"},
		Short:   "Install all the dependencies and docker images",
		Long:    `It is going to pull all the docker images necessary to spin up a Hyperledger Fabric network`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
			if err != nil {
				log.Fatalln(err)
			}

			if err := docker.CheckServerRunning(); err != nil {
				log.Fatalln(err)
			}

			config := config.GetConfig()
			for _, image := range config.FabCfg.Images {
				spinner.Spin.Start()
				spinner.Spin.Suffix = fmt.Sprintf(" Pulling %s:%s", image.Image, image.Tag)
				reader, err := cli.ImagePull(ctx, image.Image+":"+image.Tag, types.ImagePullOptions{})
				if err != nil {
					log.Fatalln(err)
				}
				if _, err := io.Copy(io.Discard, reader); err != nil {
					log.Fatalln(err)
				}
				defer reader.Close()
				spinner.Spin.Stop()
			}
		},
	}

	return cmd
}
