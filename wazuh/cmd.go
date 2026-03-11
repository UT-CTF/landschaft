package wazuh

import (
	"github.com/spf13/cobra"
)

func SetupCommand(cmd *cobra.Command) {
	setupInstallServerCmd(cmd)
}

func Run(cmd *cobra.Command) {
	installServer()
}
