package triage

import (
	"fmt"
)

func Run() {
	fmt.Println("Running triage")
	print_network_info()
	printTriageMessage()
}
