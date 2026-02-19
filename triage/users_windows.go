package triage

import (
	"fmt"
	"strings"

	"github.com/UT-CTF/landschaft/util"
)

func runUsersTriage() string {
	return parseUsersAndGroups(util.RunAndPrintScript("triage/users.ps1")) + "\t"

}

func parseUsersAndGroups(result string) string {
	lines := strings.Split(result, "\n")
	var enabled, disabled []string
	var enabledCount, disabledCount string
	current := ""
	inGroups := false

	groups := make(map[string][]string)
	var groupOrder []string
	currentGroup := ""

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "Groups:" {
			inGroups = true
			continue
		}
		if line == "" || line == "--------------------------------------------------" || line == "Users:" {
			continue
		}

		if !inGroups {
			if strings.HasPrefix(line, "Enabled Local Users") {
				current = "enabled"
				enabledCount = strings.TrimPrefix(line, "Enabled Local Users ")
				continue
			}
			if strings.HasPrefix(line, "Disabled Local Users") {
				current = "disabled"
				disabledCount = strings.TrimPrefix(line, "Disabled Local Users: ")
				continue
			}
			if current == "enabled" {
				enabled = append(enabled, line)
			} else if current == "disabled" {
				disabled = append(disabled, line)
			}
		} else {
			if strings.Contains(line, "(") && strings.Contains(line, ")") && !strings.Contains(line, "\\") {
				currentGroup = line
				groupOrder = append(groupOrder, currentGroup)
				groups[currentGroup] = []string{}
			} else if currentGroup != "" {
				groups[currentGroup] = append(groups[currentGroup], line)
			}
		}
	}

	var groupParts []string
	for _, g := range groupOrder {
		groupParts = append(groupParts, fmt.Sprintf("%s: %s", g, strings.Join(groups[g], ", ")))
	}

	return fmt.Sprintf("Enabled Local Users %s: %s; Disabled Local Users: %s: %s\t%s",
		enabledCount,
		strings.Join(enabled, ", "),
		disabledCount,
		strings.Join(disabled, ", "),
		strings.Join(groupParts, "; "),
	)
}
