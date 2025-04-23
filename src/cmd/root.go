package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "simplecloud",
	Short: "A CLI tool to manage Docker Compose apps with Traefik and Restic",
	Long:  `Install, configure, and backup self-hosted apps easily using Docker Compose, Traefik, and Restic.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
