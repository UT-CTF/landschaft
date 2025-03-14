package embed

import (
	"fmt"
	"os"
	"os/exec"
	"path"
)

func ListScripts() ([]string, error) {
	entries, err := ReadScriptDir("")
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
func ExecuteScript(scriptPath string, redirectStdout bool, additionalArgs ...string) (string, error) {
	tmpDir, err := extractEmbeddedDir(path.Dir(scriptPath))
	if err != nil {
		return "", fmt.Errorf("failed to extract embedded directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	// Prepare command execution
	fullScriptPath := path.Join(tmpDir, path.Base(scriptPath))
	cmd := exec.Command(shellName, getCommandArgs(fullScriptPath, additionalArgs...)...)
	cmd.Dir = tmpDir
	if redirectStdout {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			return "", fmt.Errorf("failed to execute script: %w", err)
		}
		return "", nil
	} else {
		output, err := cmd.CombinedOutput()
		if err != nil {
			return "", fmt.Errorf("failed to execute script: %w\nOutput: %s", err, string(output))
		}

		return string(output), nil
	}
}
