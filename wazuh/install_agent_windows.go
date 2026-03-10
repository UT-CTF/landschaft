package wazuh

import (
	"github.com/UT-CTF/landschaft/util"
)

func installAgent(agentName, managerIP, _, _, wazuhVersion string) {
	util.RunAndRedirectScript("wazuh/install_agent.ps1", managerIP, agentName, wazuhVersion)
}
