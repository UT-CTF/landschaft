package audit

import (
	"fmt"

	"github.com/UT-CTF/landschaft/util"
)

func Run() {
	fmt.Println(util.TitleColor.Render("Check SSHD"))
	checkSSHD()
}
