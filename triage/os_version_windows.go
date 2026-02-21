package triage

import (
	"fmt"
	"strings"

	"github.com/UT-CTF/landschaft/util"
)

func runOSVersionTriage() string {
	return parseOSVersion(util.RunAndPrintScript("triage/os_version.ps1"))
}

func parseOSVersion(result string) string {
	lines := strings.Split(result, "\n")
	var osName, osVersion, osConf string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "OS Name:") {
			osName = strings.TrimSpace(strings.TrimPrefix(line, "OS Name:"))
		}
		if strings.HasPrefix(line, "OS Version:") {
			osVersion = strings.TrimSpace(strings.TrimPrefix(line, "OS Version:"))
		}

		if strings.HasPrefix(line, "OS Configuration:") {
			osConf = strings.TrimSpace(strings.TrimPrefix(line, "OS Configuration:"))
		}
	}

	return fmt.Sprintf("\"%s\n%s\n\n%s\"", osName, osVersion, osConf)
}
