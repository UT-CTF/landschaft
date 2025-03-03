package cmd

import (
	"github.com/UT-CTF/landschaft/triage"

	"github.com/spf13/cobra"
)

// triageCmd represents the triage command
var triageCmd = &cobra.Command{
	Use:   "triage",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		triage.Run()
	},
}

func init() {
	rootCmd.AddCommand(triageCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// triageCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// triageCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
