package harden

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/UT-CTF/landschaft/embed"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var configFirewallArgs struct {
	outbound   bool
	outputPath string
	skipPrompt bool
}

var applyFirewallArgs struct {
	outbound   bool
	rulesFile  string
	backupPath string
}

func setupFirewallCmd(cmd *cobra.Command) {
	var firewallCmd = &cobra.Command{
		Use:   "firewall",
		Short: "Setup firewall rules",
	}

	var firewallConfigCmd = &cobra.Command{
		Use:   "config",
		Short: "Configure firewall rules to apply",
		Run: func(cmd *cobra.Command, args []string) {
			generateFirewallRules(configFirewallArgs.outbound, configFirewallArgs.outputPath, configFirewallArgs.skipPrompt)
		},
	}

	firewallConfigCmd.Flags().BoolVar(&configFirewallArgs.outbound, "outbound", false, "Configure outbound rules instead of inbound")
	firewallConfigCmd.Flags().StringVarP(&configFirewallArgs.outputPath, "output", "o", "", "Path to output the generated rules")
	firewallConfigCmd.Flags().BoolVar(&configFirewallArgs.skipPrompt, "skip", false, "Skip the confirmation prompt and select all rules")

	firewallConfigCmd.MarkFlagRequired("output")

	firewallCmd.AddCommand(firewallConfigCmd)

	var firewallApplyCmd = &cobra.Command{
		Use:   "apply",
		Short: "Apply firewall rules from a file, create a backup of existing rules, and create a scheduled task to re-apply old rules in 3 minutes",
		Run: func(cmd *cobra.Command, args []string) {
			applyFirewallRules(applyFirewallArgs.outbound, applyFirewallArgs.rulesFile, applyFirewallArgs.backupPath)
		},
	}

	firewallApplyCmd.Flags().BoolVar(&applyFirewallArgs.outbound, "outbound", false, "Apply outbound rules instead of inbound")
	firewallApplyCmd.Flags().StringVarP(&applyFirewallArgs.rulesFile, "rules", "f", "", "Path to the firewall rules file")
	firewallApplyCmd.Flags().StringVarP(&applyFirewallArgs.backupPath, "backup", "b", "firewall_backup.wfw", "Path to backup existing firewall rules")

	firewallApplyCmd.MarkFlagRequired("rules")

	firewallCmd.AddCommand(firewallApplyCmd)

	firewallFinalizeCmd := &cobra.Command{
		Use:   "finalize",
		Short: "Finalize firewall rules application by clearing the scheduled task that restores the previous firewall rules",
		Run: func(cmd *cobra.Command, args []string) {
			finalizeFirewallRules()
		},
	}

	firewallCmd.AddCommand(firewallFinalizeCmd)

	cmd.AddCommand(firewallCmd)
}

func generateFirewallRules(outbound bool, outputPath string, skipPrompt bool) {
	jsonStr := ""
	err := error(nil)

	if outbound {
		jsonStr, err = embed.ExecuteScript("harden/get_firewall_rules.ps1", false, "-Outbound")
	} else {
		jsonStr, err = embed.ExecuteScript("harden/get_firewall_rules.ps1", false)
	}

	if err != nil {
		fmt.Println("Error executing script: ", err)
		return
	}

	var raw []map[string]interface{}
	err = json.Unmarshal([]byte(jsonStr), &raw)
	if err != nil {
		fmt.Println("Error parsing JSON: ", err)
		fmt.Println("Raw JSON: ", jsonStr)
		return
	}

	var rules []map[string]string
	for _, item := range raw {
		converted := make(map[string]string)
		for k, v := range item {
			converted[k] = fmt.Sprintf("%v", v)
		}
		rules = append(rules, converted)
	}

	// any rules with enabled set to true should be pre-selected in the form
	// and then remove the enabled field for all rules since it's not needed for the actual application of the rules

	var selected []int
	for i, rule := range rules {
		if strings.ToLower(rule["Enabled"]) == "true" {
			selected = append(selected, i)
		}
		delete(rule, "Enabled")
	}

	if !skipPrompt {
		var options []huh.Option[int]
		for i, rule := range rules {
			options = append(options, huh.NewOption(rule["DisplayName"], i))
		}

		ruleSelect := huh.NewMultiSelect[int]().
			Options(options...).
			Title("Select Firewall Rules to Enable").
			Value(&selected)

		fullForm := huh.NewForm(
			huh.NewGroup(ruleSelect),
		)
		err = fullForm.Run()
		if err != nil {
			fmt.Println("Error running form: ", err)
			return
		}
	}

	var rulesToEnable []map[string]string
	for _, rule := range selected {
		rulesToEnable = append(rulesToEnable, rules[rule])
	}

	rulesToEnableJSON, err := json.MarshalIndent(rulesToEnable, "", "    ")
	if err != nil {
		fmt.Println("Error marshalling JSON: ", err)
		return
	}

	outputPath, err = filepath.Abs(outputPath)
	if err != nil {
		fmt.Println("Error getting absolute path: ", err)
		return
	}

	err = os.WriteFile(outputPath, rulesToEnableJSON, 0644)
	if err != nil {
		fmt.Println("Error writing JSON to file: ", err)
		return
	}

	fmt.Println("Rules written to: ", outputPath)
}

func applyFirewallRules(outbound bool, rulesFile string, backupPath string) {
	rulesFile, err := filepath.Abs(rulesFile)
	if err != nil {
		fmt.Println("Error getting absolute path: ", err)
		return
	}

	backupPath, err = filepath.Abs(backupPath)
	if err != nil {
		fmt.Println("Error getting absolute path: ", err)
		return
	}

	direction := "Inbound"
	if outbound {
		direction = "Outbound"
	}

	_, err = embed.ExecuteScript("harden/apply_firewall.ps1", true, "-Direction", direction, "-RulesFile", rulesFile, "-BackupFile", backupPath)
	if err != nil {
		fmt.Println("Error executing script: ", err)
		return
	}
}

func finalizeFirewallRules() {
	_, err := embed.ExecuteScript("harden/apply_firewall.ps1", true, "-ClearScheduledTask")
	if err != nil {
		fmt.Println("Error executing script: ", err)
		return
	}
}
