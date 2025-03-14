package graylog

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var genCaCmd = &cobra.Command{
	Use:   "gen-ca",
	Short: "Generate a CA for graylog",
	Run: func(cmd *cobra.Command, args []string) {
		if err := genCa(); err != nil {
			log.Error("Failed to generate CA: ", "err", err)
		}
	},
}

func setupGenCaCmd(cmd *cobra.Command) {
	cmd.AddCommand(genCaCmd)
}
