package check

import (
	"fmt"

	"github.com/UT-CTF/landschaft/util"
	"github.com/spf13/cobra"
)

// BatchCheckResult holds results of batch credential checking
type BatchCheckResult struct {
	Total        int
	Valid        int
	Invalid      int
	ValidCreds   []string
	InvalidCreds []string
}

func SetupCommand(cmd *cobra.Command) {
	setupLdapCheckCmd(cmd)
	setupKerberosCheckCmd(cmd)
	setupSmbCheckCmd(cmd)
	setupSshCheckCmd(cmd)
}

func Run(cmd *cobra.Command) {
	fmt.Println(util.ErrorStyle.Render("Error: No subcommand specified"))
	fmt.Println()
	_ = cmd.Usage()
}
