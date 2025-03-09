package embed

import (
	"embed"
	"fmt"
	"strings"
)

//go:embed windows/*/*
var scriptsFS embed.FS

var (
	scriptRootDir = "windows/"
	shellName     = "powershell"
	shellArgs     = []string{"-NoProfile", "-ExecutionPolicy", "Bypass", "-Command"}
)

func getCommandArgs(fullScriptPath string, additionalArgs ...string) []string {
	return append(shellArgs, fmt.Sprintf("%s %s", fullScriptPath, strings.Join(additionalArgs, " ")))
}
