package wazuh

import (
	"github.com/spf13/cobra"
)

var installServerCmd = &cobra.Command{
	Use:   "install-server",
	Short: "Install Wazuh manager (Linux only)",
	Run: func(cmd *cobra.Command, args []string) {
		installServer()
	},
}

func setupInstallServerCmd(cmd *cobra.Command) {
	cmd.AddCommand(installServerCmd)
}
