package harden

import (
	"fmt"

	"github.com/spf13/cobra"
)

var firewallArgs struct {
	outbound    bool
	inbound     bool
	apply       bool
	ruleFile    string
	oldRuleFile string
	backupPath  string
	removeOld   bool
	oldRulesIn  string
}

var firewallCmd = &cobra.Command{
	Use:   "firewall",
	Short: "Setup firewall rules",
	Run: func(cmd *cobra.Command, args []string) {
		// configureFirewall(firewallArgs.outbound, firewallArgs.apply, firewallArgs.ruleFile)
		if firewallArgs.apply {
			if firewallArgs.inbound == firewallArgs.outbound {
				fmt.Println("Error: You must specify either --outbound or --inbound, but not both.")
				cmd.Usage()
				return
			}

			direction := "inbound"
			if firewallArgs.outbound {
				direction = "outbound"
			}

			applyFirewallRules(firewallArgs.ruleFile, firewallArgs.backupPath, firewallArgs.oldRuleFile, direction)
		} else if firewallArgs.removeOld {
			removeOldFirewallRules(firewallArgs.oldRulesIn)
		} else {
			generateFirewallRules(firewallArgs.outbound)
		}
	},
}

func setupFirewallCmd(cmd *cobra.Command) {
	firewallCmd.Flags().BoolVar(&firewallArgs.outbound, "outbound", false, "Configure/apply outbound rules")
	firewallCmd.Flags().BoolVar(&firewallArgs.inbound, "inbound", false, "Configure/apply outbound rules")
	firewallCmd.Flags().BoolVar(&firewallArgs.apply, "apply", false, "Apply the rules")
	firewallCmd.Flags().StringVarP(&firewallArgs.ruleFile, "file", "f", "", "Path to the rules file")
	firewallCmd.Flags().StringVarP(&firewallArgs.oldRuleFile, "out", "o", "old_rules.txt", "Path to output old rules after applying new ones")
	firewallCmd.Flags().StringVarP(&firewallArgs.backupPath, "backup", "b", "firewall_backup.wfw", "Path to backup the current rules")
	firewallCmd.Flags().BoolVar(&firewallArgs.removeOld, "remove-old", false, "Remove old rules after applying new ones")
	firewallCmd.Flags().StringVarP(&firewallArgs.oldRulesIn, "old-rules", "i", "", "Path to the old rules file")

	firewallCmd.MarkFlagsRequiredTogether("file", "apply")
	firewallCmd.MarkFlagsRequiredTogether("remove-old", "old-rules")
	firewallCmd.MarkFlagsMutuallyExclusive("apply", "remove-old")
	firewallCmd.MarkFlagsMutuallyExclusive("outbound", "inbound")

	cmd.AddCommand(firewallCmd)
}
