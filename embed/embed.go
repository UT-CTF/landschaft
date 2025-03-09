package embed

import (
	"fmt"
	"os"
	"os/exec"
	"path"
)

func ListScripts() ([]string, error) {
	entries, err := readScriptDir("")
	if err != nil {
		return nil, fmt.Errorf("failed to read embedded scripts: %w", err)
	}

	var names []string
	for _, entry := range entries {
		names = append(names, entry.Name())
	}

	return names, nil
}

// ExecuteScript executes the script at the given embedded path and returns its output.
// Uses powershell on Windows and bash on other platforms.
// If bash is not available, sh is used as a fallback.
//
// Do not include the platform in the file path.
// E.g. ExecuteScript("triage/script.sh") instead of ExecuteScript("linux/triage/script.sh")
func ExecuteScript(scriptPath string) (string, error) {
	tmpDir, err := extractEmbeddedDir(path.Dir(scriptPath))
	if err != nil {
		return "", fmt.Errorf("failed to extract embedded directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	// Prepare command execution
	fullScriptPath := path.Join(tmpDir, path.Base(scriptPath))
	args := append(shellArgs, fullScriptPath)
	cmd := exec.Command(shellName, args...)
	cmd.Dir = tmpDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to execute script: %w\nOutput: %s", err, string(output))
	}

	return string(output), nil
}
