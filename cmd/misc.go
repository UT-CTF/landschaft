package cmd

import (
	"github.com/UT-CTF/landschaft/misc"
	"github.com/spf13/cobra"
)

// miscCmd represents the misc command
var miscCmd = &cobra.Command{
	Use:   "misc",
	Short: "Miscellaneous scripts that install tools or perform other tasks",
	Run: func(cmd *cobra.Command, args []string) {
		misc.Run(cmd)
	},
}

func init() {
	misc.SetupCommand(miscCmd)

	rootCmd.AddCommand(miscCmd)
}
