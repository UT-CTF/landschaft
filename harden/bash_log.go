package harden

import (
	"fmt"
	"github.com/spf13/cobra"
)

var configureBashCmd = &cobra.Command{
	Use:   "configure-bash [file path]",
	Short: "Set up bash logging",
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
		configureBash(filePath, backupPath)
	},
}

func setupConfigureBashCmd(cmd *cobra.Command) {
	configureBashCmd.Flags().StringP("file", "f", "/dev/shm/", "File path to save the new passwords or read generated passwords (csv format)")
	configureBashCmd.Flags().StringP("backup", "b", "helper", "Backup directory")

	cmd.AddCommand(configureBashCmd)
}
