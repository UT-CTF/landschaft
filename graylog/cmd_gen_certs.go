package graylog

import (
	"fmt"

	"github.com/spf13/cobra"
)

var genCertsCmd = &cobra.Command{
	Use:   "gen-certs",
	Short: "Generate certs for graylog",
	Run: func(cmd *cobra.Command, args []string) {
		if err := genCerts(); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	},
}

func setupGenCertsCmd(cmd *cobra.Command) {
	cmd.AddCommand(genCertsCmd)
}
