package triage

import (
	"fmt"

	"github.com/UT-CTF/landschaft/embed"
)

func runFirewallTriage() {
	fmt.Println("Executing triage/firewall.ps1")
	scriptOut, err := embed.ExecuteScript("triage/firewall.ps1")
	if err != nil {
		fmt.Println("Error executing script: ", err)
		return
	}
	fmt.Println(scriptOut)
}
