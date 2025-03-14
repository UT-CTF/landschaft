package harden

import (
	"fmt"

	"github.com/UT-CTF/landschaft/util"
	"github.com/spf13/cobra"
)

func SetupCommand(cmd *cobra.Command) {
	setupConfigureShellCmd(cmd)
	setupRotateLocalUsersCmd(cmd)
	setupRotatePwdCmd(cmd)
	setupFirewallCmd(cmd)
	setupAddLocalAdminCmd(cmd)
}

func Run(cmd *cobra.Command) {
	fmt.Println(util.ErrorStyle.Render("Error: No subcommand specified"))
	fmt.Println()
	_ = cmd.Usage()
}
