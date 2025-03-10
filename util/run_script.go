package util

import (
	"fmt"

	"github.com/UT-CTF/landschaft/embed"
)

func RunAndPrintScript(scriptPath string, args ...string) {
	fmt.Println("Executing " + scriptPath)
	scriptOut, err := embed.ExecuteScript(scriptPath, args...)
	if err != nil {
		fmt.Println("Error executing script: ", err)
		return
	}
	fmt.Println(scriptOut)
}
