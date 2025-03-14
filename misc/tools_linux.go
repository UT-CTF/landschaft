package misc

import (
	"github.com/UT-CTF/landschaft/util"
	"github.com/spf13/cobra"
)

func setupToolsCommand(cmd *cobra.Command) {
	toolsCmd := &cobra.Command{
		Use:   "tools",
		Short: "Install rsyslog",
		Run:   runToolsCommand,
	}

	cmd.AddCommand(toolsCmd)
}

func runToolsCommand(cmd *cobra.Command, args []string) {
	util.RunAndRedirectScript("misc/install_rsyslog.sh", args...)
}
