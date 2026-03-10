package harden

import "github.com/UT-CTF/landschaft/util"

func runHardenBaselineFirewall() {
	util.RunAndRedirectScript("harden/baseline_firewall.ps1")
}
