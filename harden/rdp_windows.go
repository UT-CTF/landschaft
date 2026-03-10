package harden

import "github.com/UT-CTF/landschaft/util"

func runHardenRDP() {
	util.RunAndRedirectScript("harden/rdp.ps1")
}
