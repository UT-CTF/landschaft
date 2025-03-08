package embed

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
)

//go:embed linux/*
var scripts embed.FS

func getScriptsOS() ([]fs.DirEntry, error) {
	return scripts.ReadDir("unix")
}

func executeScriptOS(scriptPath string, args string) (string, error) {
	// Create a temporary file to execute
	tmpFile, err := os.CreateTemp("", "script-*.sh")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	// Read the embedded script
	content, err := scripts.ReadFile("linux/" + scriptPath)
	if err != nil {
		return "", fmt.Errorf("failed to read embedded script: %w", err)
	}

	// Write to temp file and make executable
	if err := os.WriteFile(tmpFile.Name(), content, 0755); err != nil {
		return "", fmt.Errorf("failed to write script: %w", err)
	}

	cmd := exec.Command("/bin/sh", tmpFile.Name())
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to execute script: %w", err)
	}

	return string(output), nil
}
