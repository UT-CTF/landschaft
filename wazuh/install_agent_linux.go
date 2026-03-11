package wazuh

import (
	"github.com/UT-CTF/landschaft/util"
)

func installAgent(agentName, managerIP, serverUser, remoteKeyDir, _ string) {
	util.RunAndRedirectScript("wazuh/install_agent.sh", agentName, managerIP, serverUser, remoteKeyDir)
}
