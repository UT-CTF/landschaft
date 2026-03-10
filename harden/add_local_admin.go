package harden

import (
	"fmt"

	"github.com/spf13/cobra"
)

var addLocalAdminCmd = &cobra.Command{
	Use:   "add-admin [username]",
	Short: "Adds a local admin to this system",
	Long:  `Adds a local admin and prompts the user to enter the password.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		username := args[0]
		if PlanMode {
			fmt.Printf("Plan: would add local admin %q\n", username)
			return
		}
		addLocalAdmin(username)
	},
}

func setupAddLocalAdminCmd(cmd *cobra.Command) {
	cmd.AddCommand(addLocalAdminCmd)
}
