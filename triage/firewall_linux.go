package triage

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/UT-CTF/landschaft/util"
)

func runFirewallTriage() string {
	result := util.RunAndPrintScript("triage/firewall.sh") + "\t"

	if strings.Contains(result, "No supported firewall") {
		result = "No Firewall!\t"
	}

	result += getDomainInfo() + "\t"

	return result
}

func getDomainInfo() string {
	// Check 1: sssd.conf
	content, err := os.ReadFile("/etc/sssd/sssd.conf")
	if err == nil {
		fmt.Println("sssd.conf")
		return strings.TrimSpace(string(content))
	}

	// Check 2: realm list (if available)
	out, err := exec.Command("realm", "list").Output()
	if err == nil && len(out) > 0 {
		return strings.TrimSpace(string(out))
	}

	return "Not Domain Jointed"
}
