package triage

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func runOSVersionTriage() {
	file, err := os.Open("/etc/os-release")
	if err != nil {
		fmt.Println("Error reading OS release:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "=")
		if parts[0] == "NAME" {
			fmt.Println("OS:", strings.Trim(parts[1], `"`))
		} else if parts[0] == "VERSION" {
			fmt.Println("Version:", strings.Trim(parts[1], `"`))
		}
	}
}
