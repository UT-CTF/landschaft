package harden

import (
	"fmt"

	"github.com/spf13/cobra"
)

var configureShellCmd = &cobra.Command{
	Use:   "configure-shell",
	Short: "Set up shell logging",
	Long:  `Modify /etc/bash.bashrc to log all commands run by users to /var/log/commands.log`,
	Run: func(cmd *cobra.Command, args []string) {
		filePath, err := cmd.Flags().GetString("file")
		if err != nil {
			fmt.Println("Error getting file path")
			return
		}
		backupPath, err := cmd.Flags().GetString("backup")
		if err != nil {
			fmt.Println("Error getting file path")
			return
		}
		shellType, err := cmd.Flags().GetString("shell")
		if err != nil {
			fmt.Println("Error getting shell type")
			return
		}
		configureShell(filePath, backupPath, shellType)

	},
}

func setupConfigureShellCmd(cmd *cobra.Command) {
	configureShellCmd.Flags().StringP("file", "f", "/dev/shm/", "Directory path to log to")
	configureShellCmd.Flags().StringP("backup", "b", "helper", "Backup directory")
	configureShellCmd.Flags().StringP("shell", "s", "bash", "Shell type (bash or sh)")
	cmd.AddCommand(configureShellCmd)
}
