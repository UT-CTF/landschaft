package ldap

import (
	"fmt"

	"github.com/spf13/cobra"
)

func SetupCommand(cmd *cobra.Command) {
	setupGenerateLdifCmd(cmd)
	setupGeneratePasswordsCmd(cmd)
}

func Run(cmd *cobra.Command) {
	fmt.Println("Error: No subcommand specified")
	fmt.Println()
	_ = cmd.Usage()
}
