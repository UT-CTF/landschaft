package harden

import "github.com/UT-CTF/landschaft/util"

func runEnableAuditing() {
	util.RunAndRedirectScript("harden/enable_auditing.ps1")
}
