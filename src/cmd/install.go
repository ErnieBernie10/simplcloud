package cmd

import (
	"fmt"
	"os"

	"github.com/ErnieBernie10/simplecloud/src/internal"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install [app]",
	Short: "Install a Docker Compose app",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		app := args[0]
		run := internal.NewRunContext(internal.TargetDir, cmd.Context())
		err := installApp(app, run)
		if err != nil {
			fmt.Printf("Failed to install app %s: %v\n", app, err)
			os.Exit(1)
		} else {
			fmt.Printf("App %s installed successfully!\n", app)
		}
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}

func installApp(appName string, run *internal.RunContext) error {
	if !run.IsMasterRunning() {
		return fmt.Errorf("master is not running. Please start the service before installing an app")
	}

	app, err := run.GetApp(appName)
	if err != nil {
		return fmt.Errorf("failed to get app: %w", err)
	}

	if err := app.Deploy(); err != nil {
		return fmt.Errorf("failed to deploy app: %w", err)
	}

	return app.Run()
}
