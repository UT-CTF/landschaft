package wazuh

import (
	"github.com/spf13/cobra"
)

var (
	numAgents int
	agentIPs  string
)

var installServerCmd = &cobra.Command{
	Use:   "install-server",
	Short: "Install Wazuh manager and register agents (Linux only)",
	Long: `Installs the Wazuh manager from the official repository, registers the
specified agents by IP address, and saves their authentication keys to ./agent_keys/.

Keys can then be deployed to agents via the install-agent subcommand.`,
	Run: func(cmd *cobra.Command, args []string) {
		installServer(numAgents, agentIPs)
	},
}

func setupInstallServerCmd(cmd *cobra.Command) {
	installServerCmd.Flags().IntVarP(&numAgents, "num-agents", "n", 0, "Number of agents to register")
	installServerCmd.Flags().StringVarP(&agentIPs, "ips", "i", "", "Comma-separated list of agent IP addresses")
	_ = installServerCmd.MarkFlagRequired("num-agents")
	_ = installServerCmd.MarkFlagRequired("ips")
	cmd.AddCommand(installServerCmd)
}
