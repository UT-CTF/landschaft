package triage

import (
	"fmt"
	"strings"

	"github.com/UT-CTF/landschaft/util"
)

func runFirewallTriage() string {
	result := parseFirewall(util.RunAndPrintScript("triage/firewall.ps1")) + "\t"
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

		if strings.HasPrefix(trimmed, "Name") || strings.HasPrefix(trimmed, "----") || trimmed == "" || strings.HasPrefix(trimmed, "Interfaces:") {
			continue
		}

		fields := strings.Fields(trimmed)

		if strings.HasPrefix(trimmed, "Profile:") {
			parts := strings.Split(trimmed, " - ")
			profileName := strings.TrimSpace(strings.TrimPrefix(parts[0], "Profile:"))
			if len(parts) > 1 {
				profileEnabled[profileName] = strings.TrimSpace(parts[1]) == "Enabled"
			}
			continue
		}

		if strings.HasSuffix(trimmed, "Profile Settings:") {
			currentProfile = strings.TrimSuffix(trimmed, " Profile Settings:")
			continue
		}

		if currentProfile != "" {
			if strings.HasPrefix(trimmed, "State") {
				profileState[currentProfile] = strings.TrimSpace(strings.TrimPrefix(trimmed, "State"))
				continue
			}
			if strings.HasPrefix(trimmed, "Firewall Policy") {
				profilePolicy[currentProfile] = strings.TrimSpace(strings.TrimPrefix(trimmed, "Firewall Policy"))
				continue
			}
		}

		// interface line - last field is category, second to last is alias, rest is name
		if len(fields) >= 3 {
			category := fields[len(fields)-1]
			alias := fields[len(fields)-2]
			name := strings.Join(fields[:len(fields)-2], " ")
			interfaces = append(interfaces, iface{name, alias, category})
		}
	}

	var parts []string
	for _, i := range interfaces {
		enabled := "Disabled"
		check := false
		if profileEnabled[i.category] {
			enabled = "Enabled"
			check = true
		}
		state := profileState[i.category]
		policy := profilePolicy[i.category]

		if check {
			parts = append(parts, fmt.Sprintf("%s - %s \n\t%s - %s \n\tState %s \n\tFirewall Policy: %s",
				i.name, i.alias, i.category, enabled, state, policy))
		} else {
			parts = append(parts, fmt.Sprintf("%s - %s \n\t%s - %s",
				i.name, i.alias, i.category, enabled))
		}

	}

	return "\"" + strings.Join(parts, "\n\n ") + "\""

}
