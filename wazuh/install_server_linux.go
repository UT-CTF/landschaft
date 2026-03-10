package wazuh

import (
	"fmt"

	"github.com/UT-CTF/landschaft/util"
)

func installServer(numAgents int, agentIPs string) {
	util.RunAndRedirectScript("wazuh/install_server.sh", fmt.Sprintf("%d", numAgents), agentIPs)
}
