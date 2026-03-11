package wazuh

import "fmt"

func installServer(numAgents int, agentIPs string) {
	fmt.Println("Wazuh manager installation is only supported on Linux.")
	fmt.Println("Run this command on the designated Linux manager host.")
}
