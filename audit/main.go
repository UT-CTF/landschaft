package audit

import "github.com/UT-CTF/landschaft/util"

func Run() {
	util.PrintSectionTitle("Check SSHD")
	checkSSHD()
}
