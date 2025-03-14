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
		fmt.Println("Currently disabled.")
		shellConfigFile = "/etc/bash.bashrc"
		return
	} else if shellType == "sh" {
		fmt.Println("sh shell type experimental. Use at your own risk.")
		fmt.Println("Currently disabled.")
		shellConfigFile = "/etc/profile"
		return
	} else if shellType == "ssh" {
		shellConfigFile = "/etc/ssh/sshd_config"
	} else if shellType == "logger" {
		shellConfigFile = "/etc/environment"
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

	if shellType == "bash" {

		// Mainly backup, feel free to ignore .bash_history
		historyFile := absPath + "/.bash_history"
		if _, err := os.Stat(historyFile); os.IsNotExist(err) {
			file, err := os.Create(historyFile)
			if err != nil {
				fmt.Println("Error creating history file:", err)
				return
			}
			file.Close()
		}

		err = os.Chmod(historyFile, 0622)
		if err != nil {
			fmt.Println("Error setting permissions on history file:", err)
			return
		}

		logFile := absPath + "/.bash_log"
		if _, err := os.Stat(logFile); os.IsNotExist(err) {
			file, err := os.Create(logFile)
			if err != nil {
				fmt.Println("Error creating log file:", err)
				return
			}
			file.Close()
		}
		err = os.Chmod(logFile, 0622)
		if err != nil {
			fmt.Println("Error setting permissions on history file:", err)
			return
		}

		_, err = file.WriteString("\nexport HISTFILE=\"" + historyFile + "\"")
		if err != nil {
			fmt.Println("Error writing to shell config file:", err)
			return
		}
		_, err = file.WriteString("\nexport HISTTIMEFORMAT=\"%F %T \"\n")
		if err != nil {
			fmt.Println("Error writing to shell config file:", err)
			return
		}
		_, err = file.WriteString("\nexport PROMPT_COMMAND=\"history -a; echo \\\"$(date) - $(whoami) - $(pwd) - $(history 1)\\\" >> " + logFile + "\"")
		if err != nil {
			fmt.Println("Error writing to shell config file:", err)
			return
		}

		var cmd *exec.Cmd = exec.Command("/bin/bash", "-c", "shopt -s histappend")
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("Error configuring shell:", string(output))
			return
		}
	} else if shellType == "sh" {

		logFile := absPath + "/.sh_log"
		if _, err := os.Stat(logFile); os.IsNotExist(err) {
			file, err := os.Create(logFile)
			if err != nil {
				fmt.Println("Error creating log file:", err)
				return
			}
			file.Close()
		}

		err = os.Chmod(logFile, 0622)
		if err != nil {
			fmt.Println("Error setting permissions on history file:", err)
			return
		}

		/**
		_, err = file.WriteString("\ntrap 'echo \"$(date +\"%Y-%m-%d %H:%M:%S\") $(whoami) $$ $(pwd) $(history 1)\" >> " + logFile + "' DEBUG")
		if err != nil {
			fmt.Println("Error writing to shell config file:", err)
			return
		}
		*/

		// This screws up the shell, and is more noticable. Also requires /etc/profile to be used.
		// Will also log bash output

		_, err = file.WriteString("\nexport PS4='$(date +\"%Y-%m-%d %H:%M:%S\") $(whoami) $$ $(pwd) '")
		if err != nil {
			fmt.Println("Error writing to shell config file:", err)
			return
		}

		_, err = file.WriteString("\nset -x")
		if err != nil {
			fmt.Println("Error writing to shell config file:", err)
			return
		}

		_, err = file.WriteString("\nexec 2>> " + logFile)
		if err != nil {
			fmt.Println("Error writing to shell config file:", err)
			return
		}

		var cmd *exec.Cmd = exec.Command("/bin/sh", "-c", "set -o history")
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("Error configuring shell:", string(output))
			return
		}

	} else if shellType == "ssh" {
		logFile := absPath + "/.ssh_log"
		if _, err := os.Stat(logFile); os.IsNotExist(err) {
			file, err := os.Create(logFile)
			if err != nil {
				fmt.Println("Error creating log file:", err)
				return
			}
			file.Close()
		}

		err = os.Chmod(logFile, 0622)
		if err != nil {
			fmt.Println("Error setting permissions on history file:", err)
			return
		}

		_, err = file.WriteString("\nMatch all\n\tForceCommand /bin/bash -c 'echo \"$(date) - $(whoami) - $(pwd) - $SSH_CONNECTION - $(history 1)\" >> " + logFile + "; exec bash'")
		if err != nil {
			fmt.Println("Error writing to shell config file:", err)
			return
		}
	} else if shellType == "logger" {
		// Writing this in multiple places to ensure this is set. This should be relatively safe.
		_, err = file.WriteString("PROMPT_COMMAND='RETRN_VAL=$?;logger -p local6.debug \"exec_command $(whoami) [$$]: $(history 1 | sed \"s/^[ ]*[0-9]\\+[ ]*//\" )\"'")
		if err != nil {
			fmt.Println("Error writing to shell config file:", err)
		}

		shellConfigFile = "/etc/profile"
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

		_, err = file.WriteString("export PROMPT_COMMAND='logger \"User: $USER, PWD: $PWD, CMD: $(history 1 | sed \"s/^[ ]*[0-9]\\+[ ]*//\")\"'")
		if err != nil {
			fmt.Println("Error writing to shell config file:", err)
		}

		shellConfigFile = "/etc/bash.bashrc"
		err = takeBackup(shellConfigFile, backupPath)
		if err != nil {
			fmt.Println("Error taking backup:", err)
			return
		}

		file, err = os.OpenFile(shellConfigFile, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("Error opening shell config file:", err)
			return
		}
		defer file.Close()

		_, err = file.WriteString("export PROMPT_COMMAND='logger \"User: $USER, PWD: $PWD, CMD: $(history 1 | sed \"s/^[ ]*[0-9]\\+[ ]*//\")\"'")
		if err != nil {
			fmt.Println("Error writing to shell config file:", err)
		}

	}
	fmt.Println("Shell configured. Reload shell for changes to take effect.")
}
