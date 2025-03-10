package ldap

import (
	"bufio"
	"os"
	"strings"
	"testing"
)

func TestGenerateLdif(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name         string
		templateFile string
		csvFile      string
		expectedFile string
	}{
		{
			name:         "add_user",
			templateFile: "test/add_user.ldif.tpl",
			csvFile:      "test/add_user.csv",
			expectedFile: "test/add_user.ldif",
		},
		{
			name:         "add_user_to_group",
			templateFile: "test/add_user_to_group.ldif.tpl",
			csvFile:      "test/add_user_to_group.csv",
			expectedFile: "test/add_user_to_group.ldif",
		},
		{
			name:         "change_password",
			templateFile: "test/change_password.ldif.tpl",
			csvFile:      "test/change_password.csv",
			expectedFile: "test/change_password.ldif",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create temporary output file
			tmpFile, err := os.CreateTemp("", "ldif-output-*.ldif")
			if err != nil {
				t.Fatalf("Failed to create temporary file: %v", err)
			}
			defer os.Remove(tmpFile.Name())
			tmpFile.Close()

			// Generate LDIF
			err = generateLdif(tc.templateFile, tc.csvFile, tmpFile.Name())

			// Check error expectation
			if err != nil {
				t.Fatalf("Failed to generate LDIF: %v", err)
			}

			// Compare files, ignoring userPassword lines
			compareFilesIgnorePasswords(t, tc.expectedFile, tmpFile.Name())
		})
	}
}

// compareFilesIgnorePasswords compares two files line by line, ignoring lines containing "userPassword:"
// as these will have different values due to random salts used in password hashing
func compareFilesIgnorePasswords(t *testing.T, expectedPath, actualPath string) {
	expectedFile, err := os.Open(expectedPath)
	if err != nil {
		t.Fatalf("Failed to open expected file: %v", err)
	}
	defer expectedFile.Close()

	actualFile, err := os.Open(actualPath)
	if err != nil {
		t.Fatalf("Failed to open actual file: %v", err)
	}
	defer actualFile.Close()

	expectedScanner := bufio.NewScanner(expectedFile)
	actualScanner := bufio.NewScanner(actualFile)

	lineNum := 1
	for expectedScanner.Scan() && actualScanner.Scan() {
		expectedLine := expectedScanner.Text()
		actualLine := actualScanner.Text()

		// Skip comparison for userPassword lines as they contain randomly generated salts
		if strings.Contains(expectedLine, "userPassword:") && strings.Contains(actualLine, "userPassword:") {
			lineNum++
			continue
		}

		if expectedLine != actualLine {
			t.Errorf("Line %d mismatch:\nExpected: %s\nActual  : %s", lineNum, expectedLine, actualLine)
		}
		lineNum++
	}

	// Check if one file has more lines than the other
	if expectedScanner.Scan() {
		t.Errorf("Expected file has more lines than actual file")
	}
	if actualScanner.Scan() {
		t.Errorf("Actual file has more lines than expected file")
	}

	// Check for scanning errors
	if err := expectedScanner.Err(); err != nil {
		t.Errorf("Error scanning expected file: %v", err)
	}
	if err := actualScanner.Err(); err != nil {
		t.Errorf("Error scanning actual file: %v", err)
	}
}
