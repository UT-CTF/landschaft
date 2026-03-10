package harden

import "github.com/UT-CTF/landschaft/util"

func runHardenSSH() {
	util.RunAndRedirectScript("harden/ssh.sh")
}
