package cmd

import (
	"github.com/UT-CTF/landschaft/baseline"
	"github.com/spf13/cobra"
)

// baselineCmd represents the baseline command
var baselineCmd = &cobra.Command{
	Use:   "baseline",
	Short: "Takes and compares a baseline of the system",
	Run: func(cmd *cobra.Command, args []string) {
		baseline.Run(cmd)
	},
}

func init() {
	baseline.SetupCommand(baselineCmd)

	rootCmd.AddCommand(baselineCmd)
}
