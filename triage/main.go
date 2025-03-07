package triage

import (
	"fmt"
)

func Run() {
	fmt.Println("Running triage")
	printNetworkInfo()
	printTriageMessage()
}
