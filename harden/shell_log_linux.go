package harden

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func configureShell(filePath string, backupPath string, shellType string) {
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

	var shellConfigFile string
	if shellType == "bash" {
		shellConfigFile = "/etc/bash.bashrc"
	} else if shellType == "sh" {
		shellConfigFile = "/etc/profile"
	} else {
		fmt.Println("Unsupported shell type")
		return
	}

	err = takeBackup(shellConfigFile, backupPath)
	if err != nil {
		fmt.Println("Error taking backup:", err)
		return
	}

	file, err := os.OpenFile(shellConfigFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening shell config file:", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString("\nexport HISTFILE=\"" + absPath + "/.$USER\"")
	if err != nil {
		fmt.Println("Error writing to shell config file:", err)
		return
	}
	_, err = file.WriteString("\nexport HISTTIMEFORMAT=\"%F %T \"\n")
	if err != nil {
		fmt.Println("Error writing to shell config file:", err)
		return
	}
	_, err = file.WriteString("\nexport PROMPT_COMMAND=\"history -a;$PROMPT_COMMAND\"")
	if err != nil {
		fmt.Println("Error writing to shell config file:", err)
		return
	}

	var cmd *exec.Cmd
	if shellType == "bash" {
		cmd = exec.Command("/bin/bash", "-c", "shopt -s histappend")
	} else if shellType == "sh" {
		cmd = exec.Command("/bin/sh", "-c", "set -o history")
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error configuring shell:", string(output))
		return
	}

	fmt.Println("Shell configured. Reload shell for changes to take effect.")
}
