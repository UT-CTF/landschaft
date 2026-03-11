package wazuh

import (
	"github.com/UT-CTF/landschaft/util"
)

func installServer() {
	util.RunAndRedirectScript("wazuh/install_server.sh")
}
