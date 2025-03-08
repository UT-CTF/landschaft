package cmd

import (
	"github.com/UT-CTF/landschaft/harden"
	"github.com/spf13/cobra"
)

// hardenCmd represents the harden command
var hardenCmd = &cobra.Command{
	Use:   "harden",
	Short: "Various hardening scripts that need to be run individually",
	// 	Long: `A longer description that spans multiple lines and likely contains examples
	// and usage of using your command. For example:

	// Cobra is a CLI library for Go that empowers applications.
	// This application is a tool to generate the needed files
	// to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		harden.Run(cmd)
	},
}

func init() {
	harden.SetupCommand(hardenCmd)

	rootCmd.AddCommand(hardenCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// hardenCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// hardenCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
