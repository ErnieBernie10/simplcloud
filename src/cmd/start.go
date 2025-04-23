package cmd

import (
	"fmt"

	"github.com/ErnieBernie10/simplecloud/internal"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the master service",
	Long:  `Start the master service. This will start the master service and all the apps that are installed.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting master service...")
		if !internal.IsMasterRunning() {
			internal.RunMaster()
			fmt.Println("Master service started")
		} else {
			fmt.Println("Master service is already running")
		}
		// TODO: Check if master is already running
		// TODO: Check if master is installed
		// TODO: Check if master is configured
		// TODO: Check if master is up to date
		// TODO: Check if master is healthy
		// TODO: Check if master is running as root
		// TODO: Check if master is running as a service
		// TODO: Check if master is running as a container
		// TODO: Check if master is running as a systemd service
		// TODO: Check if master is running as a docker container
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
