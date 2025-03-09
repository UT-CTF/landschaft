package triage

import "github.com/UT-CTF/landschaft/util"

func Run() {
	util.PrintSectionTitle("Network")
	runNetworkTriage()
	util.PrintSectionTitle("Users & Groups")
	runUsersTriage()
	util.PrintSectionTitle("OS Version")
	runOSVersionTriage()
	util.PrintSectionTitle("Firewall")
	runFirewallTriage()
}
