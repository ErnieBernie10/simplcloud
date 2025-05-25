package cmd

import (

	"github.com/ErnieBernie10/simplecloud/internal/web"
	"github.com/spf13/cobra"
)


var serveCmd = &cobra.Command{
    Use: "serve",
	Short: "Serve web ui",
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		web.Serve()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
	
