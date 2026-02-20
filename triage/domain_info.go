package triage

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func getDomainInfo() string {
	if runtime.GOOS == "windows" {
		return getDomainInfoWindows() + "\t"
	} else {
		return getDomainInfoLinux() + "\t"
	}
}

func getDomainInfoLinux() string {
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

func getDomainInfoWindows() string {
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
