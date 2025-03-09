package embed

import (
	"embed"
	"os/exec"
)

//go:embed linux/*/*
var scriptsFS embed.FS

var (
	scriptRootDir = "linux/"
	shellName     = "bash"
	shellArgs     = []string{}
)

func init() {
	// Try to find bash in PATH
	path, err := exec.LookPath(shellName)
	if err != nil {
		// Fall back to sh if bash is not found
		if path, err = exec.LookPath("sh"); err == nil {
			shellName = path
		} else {
			panic("failed to find sh or bash in PATH")
		}
	} else {
		shellName = path
	}
}
