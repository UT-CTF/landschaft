package triage

import (
	"fmt"

	"github.com/UT-CTF/landschaft/util"
)

func Run() {
	fmt.Println(util.TitleColor.Render("Network"))
	runNetworkTriage()
	fmt.Println(util.TitleColor.Render("Users & Groups"))
	runUsersTriage()
	fmt.Println(util.TitleColor.Render("OS Version"))
	runOSVersionTriage()
	fmt.Println(util.TitleColor.Render("Firewall"))
	runFirewallTriage()
}
