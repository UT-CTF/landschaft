package triage

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func getDomainInfo() string {
	// Check 1: sssd.conf
	content, err := os.ReadFile("/etc/sssd/sssd.conf")
	if err == nil {
		fmt.Println("sssd.conf")
		return strings.TrimSpace(string(content)) + "\t"
	}

	// Check 2: realm list (if available)
	out, err := exec.Command("realm", "list").Output()
	if err == nil && len(out) > 0 {
		return strings.TrimSpace(string(out)) + "\t"
	}

	return "Not Domain Jointed" + "\t"
}
