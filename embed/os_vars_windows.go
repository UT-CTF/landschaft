package embed

import (
	"embed"
)

//go:embed windows/*/*
var scriptsFS embed.FS

var (
	scriptRootDir = "windows/"
	shellName     = "powershell"
	shellArgs     = []string{"-NoProfile", "-ExecutionPolicy", "Bypass", "-File"}
)
