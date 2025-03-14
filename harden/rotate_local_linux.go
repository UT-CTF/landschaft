package harden

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/log"
)

var defaultBlacklist = []string{"blackteam", "root", "sync"}

type user struct {
	name  string
	uid   string
	gid   string
	shell string
}

// todo: merge with triage into utils
func getLocalUsers() ([]string, error) {
	file, err := os.Open("/etc/passwd")
	if err != nil {
		return nil, fmt.Errorf("failed to open /etc/passwd: %w", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	userList := make([]user, 0)
	for scanner.Scan() {
		line := scanner.Text()
		arr := strings.Split(line, ":")
		// avoid all users that can't login
		if arr[len(arr)-1] != "/bin/false" && arr[len(arr)-1] != "/usr/sbin/nologin" {
			userList = append(userList, user{
				name:  arr[0],
				uid:   arr[2],
				gid:   arr[3],
				shell: arr[len(arr)-1],
			})
		}
	}

	// convert to []string of user names
	userNames := make([]string, 0)

upperFor:
	for _, u := range userList {
		name := u.name
		// check that name is not in blacklist
		for _, blacklisted := range defaultBlacklist {
			if name == blacklisted {
				continue upperFor
			}
		}

		userNames = append(userNames, name)
	}

	return userNames, nil
}

func applyPasswordChanges(csvPath string) {
	// Open the CSV file
	file, err := os.Open(csvPath)
	if err != nil {
		panic(fmt.Errorf("failed to open CSV file: %w", err))
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read all records
	records, err := reader.ReadAll()
	if err != nil {
		panic(fmt.Errorf("failed to read CSV file: %w", err))
	}

	// Process each record
	for _, record := range records {
		if len(record) != 2 {
			log.Error("Skipping invalid record", "record", record)
			continue
		}

		username := record[0]
		password := record[1]

		// Construct the chpasswd command
		cmd := exec.Command("chpasswd")
		cmd.Stdin = strings.NewReader(fmt.Sprintf("%s:%s", username, password))

		// Execute the command
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Error("Failed to change password for user", "user", username, "err", err, "output", string(output))
			continue
		}

		fmt.Printf("Successfully updated password for user: %s\n", username)
	}
}
