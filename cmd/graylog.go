package cmd

import (
	"github.com/UT-CTF/landschaft/graylog"
	"github.com/spf13/cobra"
)

// graylogCmd represents the harden command
var graylogCmd = &cobra.Command{
	Use:   "graylog",
	Short: "Graylog agent and server installation",
	Run: func(cmd *cobra.Command, args []string) {
		graylog.Run(cmd)
	},
}

func init() {
	graylog.SetupCommand(graylogCmd)

	rootCmd.AddCommand(graylogCmd)
}
