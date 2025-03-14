package cmd

import (
	"github.com/UT-CTF/landschaft/harden"
	"github.com/spf13/cobra"
)

// hardenCmd represents the harden command
var hardenCmd = &cobra.Command{
	Use:   "harden",
	Short: "Various hardening scripts that need to be run individually",
	Run: func(cmd *cobra.Command, args []string) {
		harden.Run(cmd)
	},
}

func init() {
	harden.SetupCommand(hardenCmd)
	rootCmd.AddCommand(hardenCmd)
}
