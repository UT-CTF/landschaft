package triage

import (
	"runtime"

	"github.com/UT-CTF/landschaft/triage/linux"
	"github.com/UT-CTF/landschaft/triage/windows"
)

func Run() {
	if runtime.GOOS == "windows" {
		windows.Run()
	} else {
		linux.Run()
	}
}
