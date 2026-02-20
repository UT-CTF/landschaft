package triage

import (
	"strings"

	"github.com/UT-CTF/landschaft/util"
)

func runFirewallTriage() string {
	result := util.RunAndPrintScript("triage/firewall.sh")
	result = strings.ReplaceAll(result, "\n", " ")
	result = strings.ReplaceAll(result, "\r", "")
	result = strings.Join(strings.Fields(result), " ")
	result = strings.ReplaceAll(result, "=", "")
	result += "\t"

	if strings.Contains(result, "No supported firewall") {
		result = "No Firewall!\t"
	}

	return result
}
