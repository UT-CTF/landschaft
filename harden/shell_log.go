package harden

import (
	"fmt"

	"github.com/spf13/cobra"
)

var configureShellCmd = &cobra.Command{
	Use:   "configure-shell",
	Short: "Set up shell logging",
	Long:  `Modify shell profile to log all commands run by users to <dir>/.bash_log or <dir>/.sh_log`,
	Run: func(cmd *cobra.Command, args []string) {
		filePath, err := cmd.Flags().GetString("directory")
		if err != nil {
			fmt.Println("Error getting directory")
			return
		}
		backupPath, err := cmd.Flags().GetString("backup")
		if err != nil {
			fmt.Println("Error getting backup")
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
	configureShellCmd.Flags().StringP("directory", "d", "/dev/shm/", "Directory path to log to")
	configureShellCmd.Flags().StringP("backup", "b", "backup", "Backup directory")
	configureShellCmd.Flags().StringP("shell", "s", "logger", "Shell type (bash, sh, ssh, logger)")
	cmd.AddCommand(configureShellCmd)
}
