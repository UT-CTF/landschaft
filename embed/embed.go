package embed

import (
	"fmt"
)

func ListScripts() ([]string, error) {
	entries, err := getScriptsOS()
	if err != nil {
		return nil, fmt.Errorf("failed to read embedded scripts: %w", err)
	}

	var names []string
	for _, entry := range entries {
		names = append(names, entry.Name())
	}

	return names, nil
}

func ExecuteScript(scriptPath string) (string, error) {
	return executeScriptOS(scriptPath)
}
