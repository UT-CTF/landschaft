package wazuh

import (
	"github.com/spf13/cobra"
)

var (
	managerIP    string
	agentName    string
	serverUser   string
	remoteKeyDir string
	wazuhVersion string
)

var installAgentCmd = &cobra.Command{
	Use:   "install-agent",
	Short: "Install Wazuh agent and register with manager",
	Long: `Installs the Wazuh agent on this host and registers it with the manager.

On Linux: fetches the pre-generated agent key from the manager via SCP, installs
the wazuh-agent package, and starts the service.

On Windows: downloads the Wazuh MSI installer, installs silently with the manager
IP embedded, and starts the WazuhSvc service.`,
	Run: func(cmd *cobra.Command, args []string) {
		installAgent(agentName, managerIP, serverUser, remoteKeyDir, wazuhVersion)
	},
}

func setupInstallAgentCmd(cmd *cobra.Command) {
	installAgentCmd.Flags().StringVar(&managerIP, "manager-ip", "", "Wazuh manager IP address")
	installAgentCmd.Flags().StringVar(&agentName, "agent-name", "", "Name to assign to this agent")
	installAgentCmd.Flags().StringVar(&serverUser, "server-user", "", "SSH username on the manager host (Linux only)")
	installAgentCmd.Flags().StringVar(&remoteKeyDir, "key-dir", "", "Remote directory containing agent_keys/ on the manager (Linux only)")
	installAgentCmd.Flags().StringVar(&wazuhVersion, "wazuh-version", "4.9.2", "Wazuh agent version to install (Windows only)")
	_ = installAgentCmd.MarkFlagRequired("manager-ip")
	_ = installAgentCmd.MarkFlagRequired("agent-name")
	cmd.AddCommand(installAgentCmd)
}
