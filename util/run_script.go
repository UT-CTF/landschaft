package util

import (
	"fmt"

	"github.com/UT-CTF/landschaft/embed"
	"github.com/charmbracelet/log"
)

func RunAndPrintScript(scriptPath string, args ...string) {
	fmt.Println("Executing " + scriptPath)
	scriptOut, err := embed.ExecuteScript(scriptPath, false, args...)
	if err != nil {
		log.Error("Failed to execute script", "script", scriptPath, "err", err)
		return
	}
	fmt.Println(scriptOut)
}

func RunAndRedirectScript(scriptPath string, args ...string) {
	fmt.Println("Executing " + scriptPath)
	_, err := embed.ExecuteScript(scriptPath, true, args...)
	if err != nil {
		fmt.Println("Error executing script: ", err)
		return
	}
}
