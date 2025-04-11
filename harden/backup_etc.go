package harden

import (
	"fmt"

	"github.com/spf13/cobra"
)

var backupEtcCmd = &cobra.Command{
	Use:   "backup-etc",
	Short: "backup the /etc/ folder",
	Long:  `adds a backdoor`,
	Run: func(cmd *cobra.Command, args []string) {
		backupPath, err := cmd.Flags().GetString("backup-directory")
		if err != nil {
			fmt.Println("Error getting directory")
			return
		}
		backup_etc(backupPath)

	},
}

var restoreEtcCmd = &cobra.Command{
	Use:   "restore-etc",
	Short: "restores the /etc/ folder",
	Long:  `restores the /etc/ folder from the provided tar file`,
	Run: func(cmd *cobra.Command, args []string) {
		backupPath, err := cmd.Flags().GetString("restore-directory")
		if err != nil {
			fmt.Println("Error getting directory")
			return
		}
		restore_etc(backupPath)

	},
}

func setupBackupEtcCmd(cmd *cobra.Command) {
	backupEtcCmd.Flags().StringP("backup-directory", "d", "/", "Directory path to backup to")
	cmd.AddCommand(backupEtcCmd)
}

func setupRestoreEtcCmd(cmd *cobra.Command) {
	restoreEtcCmd.Flags().StringP("restore-directory", "d", "/", "Directory path to backup to")
	cmd.AddCommand(backupEtcCmd)
}
