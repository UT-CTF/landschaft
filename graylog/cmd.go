package graylog

import (
	"fmt"

	"github.com/spf13/cobra"
)

func SetupCommand(cmd *cobra.Command) {
	setupGenCertsCmd(cmd)
}

func Run(cmd *cobra.Command) {
	fmt.Println("Error: No subcommand specified")
	fmt.Println()
	_ = cmd.Usage()
}
