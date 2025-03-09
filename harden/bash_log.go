package harden

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

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

func configureBash(filePath string, backupPath string) {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		fmt.Println("Error getting absolute path")
		return
	}


	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		err = os.Mkdir(absPath, 0755)
		if err != nil {
			fmt.Println("Error creating history directory:", err)
			return
		}
	}

	err = takeBackup("/etc/bash.bashrc", "helper")
	if err != nil {
		fmt.Println("Error taking backup:", err)
		return
	}

	file, err := os.OpenFile("/etc/bash.bashrc", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening bash.bashrc:", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString("\nexport HISTFILE=\"" + absPath + ".$USER\"")
	if err != nil {
		fmt.Println("Error writing to bash.bashrc:", err)
		return
	}
	_, err = file.WriteString("\nexport HISTTIMEFORMAT=\"%F %T \"\n")
	if err != nil {
		fmt.Println("Error writing to bash.bashrc:", err)
		return
	}
	_, err = file.WriteString("\nexport PROMPT_COMMAND=\"history -a;$PROMPT_COMMAND\"")
	if err != nil {
		fmt.Println("Error writing to bash.bashrc:", err)
		return
	}

	cmd := exec.Command("/bin/bash", "-c", "shopt -s histappend")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error configuring bash:", string(output))
		return
	}

	fmt.Println("Bash configured. Reload shell for changes to take effect.")
}

func setupConfigureBashCmd(cmd *cobra.Command) {
	configureBashCmd.Flags().StringP("file", "f", "/dev/shm/", "File path to save the new passwords or read generated passwords (csv format)")
	configureBashCmd.Flags().StringP("backup", "b", "helper", "Backup directory")

	cmd.AddCommand(configureBashCmd)
}
