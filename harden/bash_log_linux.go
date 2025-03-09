package harden

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)



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

	err = takeBackup("/etc/bash.bashrc", backupPath)
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

