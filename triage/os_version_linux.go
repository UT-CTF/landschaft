package triage

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func runOSVersionTriage() string {
	file, err := os.Open("/etc/os-release")
	if err != nil {
		fmt.Println("Error reading OS release:", err)
		return "err"
	}
	defer file.Close()

	var result string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "=")
		if parts[0] == "NAME" {
			result += strings.Trim(parts[1], `"`) + "\n"
			fmt.Println("OS:", strings.Trim(parts[1], `"`))
		} else if parts[0] == "VERSION" {
			fmt.Println("Version:", strings.Trim(parts[1], `"`))
			result += strings.Trim(parts[1], `"`)
		}
	}
	return "\"" + result + "\""
}
