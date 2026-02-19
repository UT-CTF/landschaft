package triage

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/UT-CTF/landschaft/util"
)

func runFirewallTriage() string {
	result := parseFirewall(util.RunAndPrintScript("triage/firewall.ps1")) + "\t"
	result += getDomainStatus() + "\t"
	return result
}

func parseFirewall(result string) string {
	lines := strings.Split(result, "\n")

	type iface struct {
		name     string
		alias    string
		category string
	}

	var interfaces []iface
	profileEnabled := make(map[string]bool)
	profileState := make(map[string]string)
	profilePolicy := make(map[string]string)
	currentProfile := ""

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// skip header lines
		if strings.HasPrefix(trimmed, "Name") || strings.HasPrefix(trimmed, "----") || trimmed == "" {
			continue
		}

		// interface lines - 3 fields
		fields := strings.Fields(trimmed)
		if len(fields) == 3 && !strings.HasPrefix(trimmed, "Profile:") && !strings.HasSuffix(trimmed, "Settings:") && !strings.HasPrefix(trimmed, "State") && !strings.HasPrefix(trimmed, "Firewall") {
			interfaces = append(interfaces, iface{fields[0], fields[1], fields[2]})
			continue
		}

		// Profile: X - Enabled/Disabled
		if strings.HasPrefix(trimmed, "Profile:") {
			parts := strings.Split(trimmed, " - ")
			profileName := strings.TrimSpace(strings.TrimPrefix(parts[0], "Profile:"))
			if len(parts) > 1 {
				profileEnabled[profileName] = strings.TrimSpace(parts[1]) == "Enabled"
			}
			continue
		}

		// X Profile Settings:
		if strings.HasSuffix(trimmed, "Profile Settings:") {
			currentProfile = strings.TrimSuffix(trimmed, " Profile Settings:")
			continue
		}

		if currentProfile != "" {
			if strings.HasPrefix(trimmed, "State") {
				profileState[currentProfile] = strings.TrimSpace(strings.TrimPrefix(trimmed, "State"))
			}
			if strings.HasPrefix(trimmed, "Firewall Policy") {
				profilePolicy[currentProfile] = strings.TrimSpace(strings.TrimPrefix(trimmed, "Firewall Policy"))
			}
		}
	}

	var parts []string
	for _, i := range interfaces {
		enabled := "Disabled"
		if profileEnabled[i.category] {
			enabled = "Enabled"
		}
		state := profileState[i.category]
		policy := profilePolicy[i.category]
		parts = append(parts, fmt.Sprintf("%s - %s - %s (%s; State %s; Firewall Policy: %s)",
			i.name, i.alias, i.category, enabled, state, policy))
	}

	return strings.Join(parts, "; ")
}

func getDomainStatus() string {
	out, err := exec.Command("wmic", "computersystem", "get", "domain").Output()
	if err != nil {
		return "Not Domain Joined"
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || line == "Domain" {
			continue
		}
		if line == "WORKGROUP" {
			return "Not Domain Joined"
		}
		return line
	}

	return "Not Domain Joined"
}
