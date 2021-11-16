package network

import (
	"log"

	"github.com/spf13/cobra"
)

// NewCmdNetwork is the root command for all operations related to a network
func NewCmdNetwork() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "network",
		Short: "All operations to install, run and stop a Hyperledger Fabric network",
		Long:  `All operations to install, run and stop a Hyperledger Fabric network.`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Help(); err != nil {
				log.Fatalln(err)
			}
		},
	}

	cmd.AddCommand(NewCmdNetworkInstall())

	return cmd
}
