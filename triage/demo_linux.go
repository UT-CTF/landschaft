package triage

import (
	"fmt"

	"github.com/UT-CTF/landschaft/embed"
)

func printTriageMessage() {
	fmt.Println("Running linux triage demo")

	scriptOutput, err := embed.ExecuteScript("test.sh")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(scriptOutput)
}
