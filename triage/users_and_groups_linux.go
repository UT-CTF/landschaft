package triage

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func printUsers() {
	file, err := os.Open("/etc/passwd")
	if err != nil {
		fmt.Println("Error reading /etc/passwd:", err)
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	fmt.Println("Listing users:")
	fmt.Printf("%-32s|%-4s|%-4s|%s\n", "username", "UID", "GID", "shell")
	fmt.Println("------------------------------------------------")
	for scanner.Scan() {
		line := scanner.Text()
		arr := strings.Split(line, ":")
		// avoid all users that can't login
		if arr[len(arr)-1] != "/bin/false" && arr[len(arr)-1] != "/usr/sbin/nologin" {
			fmt.Printf("%-32s|%-4s|%-4s|%s\n", arr[0], arr[2], arr[2], arr[len(arr)-1])
		}
	}
}

func printGroups() {
	file, err := os.Open("/etc/group")
	if err != nil {
		fmt.Println("Error reading /etc/group:", err)
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	fmt.Println("Listing groups:")
	fmt.Printf("%-32s|%-4s|%s\n", "groupname", "GID", "group list")
	fmt.Println("------------------------------------------------")
	for scanner.Scan() {
		line := scanner.Text()
		arr := strings.Split(line, ":")
		// avoid all groups that no one is a part of
		if arr[3] != "" {
			fmt.Printf("%-32s|%-4s|%s\n", arr[0], arr[2], arr[3])
		}
	}
}
