package harden

import (
	"github.com/spf13/cobra"
)

var addLocalAdminCmd = &cobra.Command{
	Use:   "add-admin [username]",
	Short: "Adds a local admin to this system",
	Long:  `Adds a local admin and prompts the user to enter the password.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		username := args[0]
		addLocalAdmin(username)
	},
}

func setupAddLocalAdminCmd(cmd *cobra.Command) {
	cmd.AddCommand(addLocalAdminCmd)
}
