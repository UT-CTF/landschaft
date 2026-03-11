package cmd

import (
	"github.com/UT-CTF/landschaft/wazuh"
	"github.com/spf13/cobra"
)

// wazuhCmd represents the wazuh command
var wazuhCmd = &cobra.Command{
	Use:   "wazuh",
	Short: "Wazuh agent and server installation",
	Run: func(cmd *cobra.Command, args []string) {
		wazuh.Run(cmd)
	},
}

func init() {
	wazuh.SetupCommand(wazuhCmd)
	rootCmd.AddCommand(wazuhCmd)
}
