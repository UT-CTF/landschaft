package triage

import (
	"fmt"

	"github.com/UT-CTF/landschaft/embed"
)

func runOSVersionTriage() {
	fmt.Println("Executing triage/os_version.ps1")
	scriptOut, err := embed.ExecuteScript("triage/os_version.ps1")
	if err != nil {
		fmt.Println("Error executing script: ", err)
		return
	}
	fmt.Println(scriptOut)
}
