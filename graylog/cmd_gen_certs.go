package graylog

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var genCertsCmd = &cobra.Command{
	Use:   "gen-certs",
	Short: "Generate certs for graylog",
	Run: func(cmd *cobra.Command, args []string) {
		if err := genCerts(); err != nil {
			log.Error("Failed to generate cert:", "err", err)
		}
	},
}

func setupGenCertsCmd(cmd *cobra.Command) {
	cmd.AddCommand(genCertsCmd)
}
