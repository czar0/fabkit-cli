package network

import (
	"context"
	"fmt"
	"log"

	"github.com/czar0/fabkit-cli/internal/docker"
	"github.com/czar0/fabkit-cli/internal/spinner"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

// TODO: Image names should set in some const or config file, while the version should be set at runtime (replacing default value in case the -v flag is a valid Fabric version)
var (
	images = []string{"hyperledger/fabric-orderer:2.3", "hyperledger/fabric-peer:2.3", "hyperledger/fabric-ccenv:2.3", "hyperledger/fabric-ca:1.5", "hyperledger/fabric-couchdb:0.4.22"}
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

			for _, image := range images {
				spinner.Spin.Start()
				spinner.Spin.Suffix = fmt.Sprintf(" Pulling %s", image)
				reader, err := cli.ImagePull(ctx, image, types.ImagePullOptions{})
				if err != nil {
					log.Fatalln(err)
				}
				defer reader.Close()
				spinner.Spin.Stop()
			}
		},
	}

	return cmd
}
