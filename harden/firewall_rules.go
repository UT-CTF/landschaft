package harden

import "github.com/spf13/cobra"

var firewallCmd = &cobra.Command{
	Use:   "firewall",
	Short: "Setup firewall rules",
	Run: func(cmd *cobra.Command, args []string) {
		configureFirewall()
	},
}

func setupFirewallCmd(cmd *cobra.Command) {
	cmd.AddCommand(firewallCmd)
}
