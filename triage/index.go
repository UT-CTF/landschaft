package triage

import (
	"fmt"
	"os"

	"github.com/UT-CTF/landschaft/triage/linux"
	"github.com/UT-CTF/landschaft/triage/windows"
	"github.com/UT-CTF/landschaft/util"
)

func Run() {
	if util.IsWindows() {
		windows.Run()
	} else if util.IsLinux() {
		linux.Run()
	} else {
		fmt.Println("Unsupported OS")
		os.Exit(1)
	}
}
