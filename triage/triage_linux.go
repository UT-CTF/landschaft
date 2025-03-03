package triage

import (
	"fmt"

	"github.com/UT-CTF/landschaft/triage/linux"
)

func Run() {
	fmt.Println("Running general linux triage")
	linux.Run()
}
