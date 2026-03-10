package harden

import (
	"fmt"
	"os"
	"path/filepath"
)

func configureShell(filePath string, backupPath string, shellType string) {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		fmt.Println("Error getting absolute path")
		return
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		if err = os.Mkdir(absPath, 0755); err != nil {
			fmt.Println("Error creating history directory:", err)
			return
		}
	}

	var shellConfigFile string
	switch shellType {
	case "bash":
		fmt.Println("Currently disabled.")
		return
	case "sh":
		fmt.Println("sh shell type experimental. Use at your own risk.")
		fmt.Println("Currently disabled.")
		return
	case "ssh":
		shellConfigFile = "/etc/ssh/sshd_config"
	case "logger":
		shellConfigFile = "/etc/environment"
	default:
		fmt.Println("Unsupported shell type")
		return
	}

	if err = takeBackup(shellConfigFile, backupPath); err != nil {
		fmt.Println("Error taking backup:", err)
		return
	}

	file, err := os.OpenFile(shellConfigFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening shell config file:", err)
		return
	}
	defer file.Close()

	switch shellType {
	case "ssh":
		logFile := absPath + "/.ssh_log"
		if _, err := os.Stat(logFile); os.IsNotExist(err) {
			lf, err := os.Create(logFile)
			if err != nil {
				fmt.Println("Error creating log file:", err)
				return
			}
			lf.Close()
		}
		if err = os.Chmod(logFile, 0622); err != nil {
			fmt.Println("Error setting permissions on log file:", err)
			return
		}
		line := "\nMatch all\n\tForceCommand /bin/bash -c 'echo \"$(date) - $(whoami) - $(pwd) - $SSH_CONNECTION - $(history 1)\" >> " + logFile + "; exec bash'"
		if _, err = file.WriteString(line); err != nil {
			fmt.Println("Error writing to shell config file:", err)
			return
		}

	case "logger":
		loggerLine := "PROMPT_COMMAND='RETRN_VAL=$?;logger -p local6.debug \"exec_command $(whoami) [$$]: $(history 1 | sed \"s/^[ ]*[0-9]\\+[ ]*//\" )\"'"
		if _, err = file.WriteString(loggerLine); err != nil {
			fmt.Println("Error writing to shell config file:", err)
		}
		for _, cfg := range []string{"/etc/profile", "/etc/bash.bashrc"} {
			if err = takeBackup(cfg, backupPath); err != nil {
				fmt.Println("Error taking backup:", err)
				return
			}
			f, err := os.OpenFile(cfg, os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Println("Error opening shell config file:", err)
				return
			}
			defer f.Close()
			exportLine := "export PROMPT_COMMAND='logger \"User: $USER, PWD: $PWD, CMD: $(history 1 | sed \"s/^[ ]*[0-9]\\+[ ]*//\")\"'"
			if _, err = f.WriteString(exportLine); err != nil {
				fmt.Println("Error writing to shell config file:", err)
			}
		}
	}

	fmt.Println("Shell configured. Reload shell for changes to take effect.")
}
