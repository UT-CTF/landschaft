package triage

import (
	"fmt"

	"github.com/UT-CTF/landschaft/embed"
)

func runUsersTriage() {
	fmt.Println("Executing triage/users.ps1")
	scriptOut, err := embed.ExecuteScript("triage/users.ps1")
	if err != nil {
		fmt.Println("Error executing script: ", err)
		return
	}
	fmt.Println(scriptOut)
}
