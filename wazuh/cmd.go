package wazuh

import (
	"fmt"

	"github.com/UT-CTF/landschaft/util"
	"github.com/spf13/cobra"
)

func SetupCommand(cmd *cobra.Command) {
	setupInstallServerCmd(cmd)
	setupInstallAgentCmd(cmd)
}

func Run(cmd *cobra.Command) {
	fmt.Println(util.ErrorStyle.Render("Error: No subcommand specified"))
	fmt.Println()
	_ = cmd.Usage()
}
