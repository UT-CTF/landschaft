package harden

import (
	"fmt"

	"github.com/spf13/cobra"
)

func setupCCDCCommands(cmd *cobra.Command) {
	sshCmd := &cobra.Command{
		Use:   "ssh",
		Short: "Harden SSH server (Linux): ensure secure sshd settings",
		Run: func(c *cobra.Command, args []string) {
			if PlanMode {
				fmt.Println("Plan: would ensure sshd_config has PermitRootLogin no and secure ciphers.")
				return
			}
			runHardenSSH()
		},
	}
	rdpCmd := &cobra.Command{
		Use:   "rdp",
		Short: "Harden RDP (Windows): NLA and restrict settings",
		Run: func(c *cobra.Command, args []string) {
			if PlanMode {
				fmt.Println("Plan: would enable NLA and restrict RDP settings.")
				return
			}
			runHardenRDP()
		},
	}
	lockAccountsCmd := &cobra.Command{
		Use:   "lock-accounts",
		Short: "Disable guest and high-risk default accounts",
		Run: func(c *cobra.Command, args []string) {
			if PlanMode {
				fmt.Println("Plan: would disable guest/test/default accounts.")
				return
			}
			runHardenLockAccounts()
		},
	}
	baselineFirewallCmd := &cobra.Command{
		Use:   "baseline-firewall",
		Short: "Apply conservative baseline firewall (allow established, common ports)",
		Run: func(c *cobra.Command, args []string) {
			if PlanMode {
				fmt.Println("Plan: would apply baseline firewall rules (allow established + common scored ports).")
				return
			}
			runHardenBaselineFirewall()
		},
	}
	cmd.AddCommand(sshCmd)
	cmd.AddCommand(rdpCmd)
	cmd.AddCommand(lockAccountsCmd)
	cmd.AddCommand(baselineFirewallCmd)
}
