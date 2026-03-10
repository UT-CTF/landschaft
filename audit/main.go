package audit

import (
	"fmt"
	"time"

	"github.com/UT-CTF/landschaft/util"
)

type Options struct {
	CheckSSHD     bool
	CheckVersions bool

	// MaxPackages caps the number of packages queried against OSV in a single run.
	// 0 means "use default".
	MaxPackages int

	// Timeout applies to OSV network calls.
	// 0 means "use default".
	Timeout time.Duration
}

func DefaultOptions() Options {
	return Options{
		CheckSSHD:     true,
		CheckVersions: true,
		MaxPackages:   50,
		Timeout:       5 * time.Second,
	}
}

func Run(opts Options) {
	if opts.CheckSSHD {
		fmt.Println(util.TitleColor.Render("Check SSHD"))
		checkSSHD()
	}

	if opts.CheckVersions {
		fmt.Println()
		fmt.Println(util.TitleColor.Render("Software versions (OSV)"))
		checkSoftwareVersions(opts)
	}
}
