package embed

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
)

//go:embed windows/*
var scripts embed.FS

func getScriptsOS() ([]fs.DirEntry, error) {
	return scripts.ReadDir("windows")
}

func executeScriptOS(scriptPath string) (string, error) {
	// Create a temporary file to execute
	tmpFile, err := os.CreateTemp("", "script-*.ps1")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	// Read the embedded script
	content, err := scripts.ReadFile("windows/" + scriptPath)
	if err != nil {
		return "", fmt.Errorf("failed to read embedded script: %w", err)
	}

	// Write to temp file
	if err := os.WriteFile(tmpFile.Name(), content, 0644); err != nil {
		return "", fmt.Errorf("failed to write script: %w", err)
	}
	tmpFile.Close()

	cmd := exec.Command("powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-File", tmpFile.Name())
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to execute script: %w", err)
	}

	return string(output), nil
}
