package ldap

import (
	"encoding/csv"
	"fmt"
	"os"
	"path"
	"slices"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/charmbracelet/x/term"
	"github.com/jsimonetti/pwscheme/ssha"
)

var funcMap = template.FuncMap{
	"fromCsv":            fromCsv,
	"encodeLdifPassword": encodeLdifPassword,
}

var csvHeader []string

type Data struct {
	CsvRows [][]string
}

type UserPass struct {
	Username string
	Password string
}

func generateLdif(templatePath string, csvPath string, outputPath string) error {
	// Open CSV file
	csvData, err := os.Open(csvPath)
	if err != nil {
		return fmt.Errorf("failed to open CSV file: %v", err)
	}
	defer csvData.Close()
	csvReader := csv.NewReader(csvData)
	csvHeader, err = csvReader.Read()
	if err != nil {
		return fmt.Errorf("failed to read csv header: %v", err)
	}

	csvRows, err := csvReader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read csv rows: %v", err)
	}

	// Create output file
	output, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer output.Close()

	// Prepare template
	t := template.New(path.Base(templatePath)).Funcs(funcMap).Funcs(sprig.FuncMap())
	t, err = t.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("failed to parse template file: %v", err)
	}

	err = t.Execute(output, Data{
		CsvRows: csvRows,
	})
	if err != nil {
		return fmt.Errorf("failed to execute template: %v", err)

	}

	return nil
}

func generateLdifSingle(templatePath string, outputPath string) error {

	// Prompt for username and password
	var username, password string
	fmt.Print("Enter username: ")
	fmt.Scanln(&username)

	var match bool
	for !match {
		fmt.Print("Enter password: ")
		passwordBytes, err := term.ReadPassword(os.Stdin.Fd())
		if err != nil {
			return fmt.Errorf("failed to read password: %v", err)
		}
		password = string(passwordBytes)
		fmt.Println()

		fmt.Print("Confirm password: ")
		var passwordConfirm string
		passwordConfBytes, err := term.ReadPassword(os.Stdin.Fd())
		if err != nil {
			return fmt.Errorf("failed to read password: %v", err)
		}
		passwordConfirm = string(passwordConfBytes)

		if password == passwordConfirm {
			match = true
		} else {
			fmt.Println("Passwords do not match. Please try again.")
		}
	}

	// Create output file
	output, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer output.Close()

	// Prepare template
	t := template.New(path.Base(templatePath)).Funcs(funcMap).Funcs(sprig.FuncMap())
	t, err = t.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("failed to parse template file: %v", err)
	}

	err = t.Execute(output, UserPass{
		Username: username,
		Password: password,
	})
	if err != nil {
		return fmt.Errorf("failed to execute template: %v", err)

	}

	fmt.Printf("LDIF file generated at %s\n", outputPath)

	return nil
}

func fromCsv(colName string, row []string) (string, error) {
	if csvHeader == nil {
		return "", fmt.Errorf("csv header not set")
	}

	idx := slices.Index(csvHeader, colName)
	if idx == -1 {
		return "", fmt.Errorf("column %s not found in csv header", colName)
	}

	return row[idx], nil
}

func encodeLdifPassword(password string) string {
	encodedPassword, err := ssha.Generate(password, 20)
	if err != nil {
		fmt.Printf("failed to encode password: %v", err)
		return ""
	}

	return encodedPassword
}
