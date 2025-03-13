package ldap

import (
	"encoding/csv"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/UT-CTF/landschaft/util"
)

// FindAllUsers finds all LDAP users using ldapsearch
func FindAllUsers(baseDn string, extraArgs []string, excludedUsers []string) ([]string, error) {
	// Start with default arguments
	args := []string{"-x", "-LLL", "-b", baseDn, "(objectclass=posixAccount)"}

	// Append extra arguments if provided
	args = append(args, extraArgs...)

	searchCmd := exec.Command("ldapsearch", args...)
	output, err := searchCmd.CombinedOutput()
	if err != nil {
		// Include the command output in the error message
		return nil, fmt.Errorf("failed to execute ldapsearch: %v\nCommand output:\n%s", err, output)
	}

	// Always exclude these users
	alwaysExcluded := map[string]bool{
		"blackteam": true,
		"root":      true,
	}

	// Add user-provided exclusions
	userExcluded := make(map[string]bool)
	for _, user := range excludedUsers {
		userExcluded[user] = true
	}

	var users []string
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "dn: uid=") {
			// Extract the username between "uid=" and the comma
			uidPart := strings.Split(line, ",")[0]
			username := strings.TrimPrefix(uidPart, "dn: uid=")

			// Skip excluded users
			if alwaysExcluded[username] || userExcluded[username] {
				continue
			}

			users = append(users, username)
		}
	}

	return users, nil
}

// GeneratePasswordsCSV generates random passwords for all users and writes to a CSV file
func GeneratePasswordsCSV(baseDn, outputPath string, passwordLength uint, allowedChars string, extraArgs []string, excludedUsers []string) error {
	// Check if file already exists
	if _, err := os.Stat(outputPath); err == nil {
		return fmt.Errorf("output file %s already exists, refusing to overwrite", outputPath)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to check if output file exists: %v", err)
	}

	// Create the directory for the output file if it doesn't exist
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create directory %s: %v", dir, err)
	}

	// Find all users
	users, err := FindAllUsers(baseDn, extraArgs, excludedUsers)
	if err != nil {
		return err
	}

	// Create output file
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer file.Close()

	// Write CSV header and user entries
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	if err := writer.Write([]string{"Username", "Password"}); err != nil {
		return fmt.Errorf("failed to write CSV header: %v", err)
	}

	// Write rows for each user with a new random password
	for _, user := range users {
		password := util.GenerateRandomPassword(passwordLength, allowedChars, true)
		if err := writer.Write([]string{user, password}); err != nil {
			return fmt.Errorf("failed to write CSV row: %v", err)
		}
	}

	return nil
}
