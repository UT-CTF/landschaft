package triage

import (
	"os/exec"
	"strings"
)

func getDomainInfo() string {
	out, err := exec.Command("wmic", "computersystem", "get", "domain").Output()
	if err != nil {
		return "Not Domain Joined" + "\t"
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || line == "Domain" {
			continue
		}
		if line == "WORKGROUP" {
			return "Not Domain Joined" + "\t"
		}
		return line + "\t"
	}

	return "Not Domain Joined" + "\t"
}
