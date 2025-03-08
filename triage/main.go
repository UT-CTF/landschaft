package triage

func Run() {
	// fmt.Println("Running triage")
	// printNetworkInfo()
	runUsersTriage()
	runOSVersionTriage()
	runFirewallTriage()
	// printTriageMessage()
}
