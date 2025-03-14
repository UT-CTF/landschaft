package cmd

import (
	"github.com/UT-CTF/landschaft/serve"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve [directory]",
	Short: "Serve directory over https",
	Run: func(cmd *cobra.Command, args []string) {
		serve.Run(cmd, args)
	},
}

func init() {
	serve.SetupCommand(serveCmd)

	rootCmd.AddCommand(serveCmd)
}
