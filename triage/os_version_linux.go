package triage

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func printOSVersion() {
	file, err := os.Open("/etc/os-release")
	if err != nil {
		fmt.Println("Error reading OS release:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "PRETTY_NAME=") {
			fmt.Println("OS:", strings.Trim(line[12:], `"`))
		} else if strings.HasPrefix(line, "VERSION=") {
			fmt.Println("Version:", strings.Trim(line[9:], `"`))
		}
	}
}
