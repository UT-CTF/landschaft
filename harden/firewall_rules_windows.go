package harden

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/UT-CTF/landschaft/embed"
	"github.com/UT-CTF/landschaft/util"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

const OutboundRulesPath = "firewall_rules_outbound.json"
const InboundRulesPath = "firewall_rules_inbound.json"

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

func generateFirewallRules(outbound bool) {
	var rulesPath = InboundRulesPath
	if outbound {
		rulesPath = OutboundRulesPath
	}

	jsonStr, err := embed.ExecuteScript("harden/firewall.ps1", false, "-RulePath", rulesPath)
	if err != nil {
		fmt.Println("Error executing script: ", err)
		return
	}

	var rules []map[string]string
	err = json.Unmarshal([]byte(jsonStr), &rules)
	if err != nil {
		fmt.Println("Error parsing JSON: ", err)
		return
	}

	sort.Slice(rules, func(i, j int) bool {
		return rules[i]["DisplayName"] < rules[j]["DisplayName"]
	})

	var selected []int

	ruleSelect := huh.NewMultiSelect[int]().OptionsFunc(func() []huh.Option[int] {
		var ruleNames []huh.Option[int]
		for i, rule := range rules {
			ruleNames = append(ruleNames, huh.NewOption(rule["DisplayName"], i))
		}
		return ruleNames
	},
		nil,
	).Title("Select Firewall Rules to Enable").
		Value(&selected)

	var outputPath string
	defaultPath := rulesPath
	pathInput := huh.NewInput().Title("Rules file path").Placeholder(defaultPath).Value(&outputPath)

	fullForm := huh.NewForm(
		huh.NewGroup(ruleSelect, pathInput),
	)
	err = fullForm.Run()
	if err != nil {
		fmt.Println("Error running form: ", err)
		return
	}

	outputPath = strings.TrimSpace(outputPath)
	if len(outputPath) == 0 {
		outputPath = defaultPath
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

func applyFirewallRules(rulesFile string, backupPath string, oldRulesFile string, direction string) {
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

	oldRulesFile, err = filepath.Abs(oldRulesFile)
	if err != nil {
		fmt.Println("Error getting absolute path: ", err)
		return
	}

	fmt.Println("This will save IDs for all existing rules to a file and apply the selected rules.")
	fmt.Println("This will also create a backup of the current rules.")
	fmt.Println("The new rules are in the file: ", rulesFile)
	fmt.Print("Are you sure you want to continue? (y/n) ")
	var confirm string
	fmt.Scanln(&confirm)
	if confirm != "y" {
		fmt.Println("Aborting.")
		return
	}

	util.RunAndPrintScript("harden/firewall.ps1", "-RulePath", "'"+rulesFile+"'", "-BackupPath", "'"+backupPath+"'", "-OldRulesPath", "'"+oldRulesFile+"'", "-Direction", "'"+direction+"'", "-Apply")
}

func removeOldFirewallRules(oldRulesFile string) {
	oldRulesFile, err := filepath.Abs(oldRulesFile)
	if err != nil {
		fmt.Println("Error getting absolute path: ", err)
		return
	}

	fmt.Println("This will remove ALL the old rules from the system.")
	fmt.Print("Are you sure you want to continue? (y/n) ")
	var confirm string
	fmt.Scanln(&confirm)
	if confirm != "y" {
		fmt.Println("Aborting.")
		return
	}

	util.RunAndPrintScript("harden/firewall.ps1", "-OldRulesPath", "'"+oldRulesFile+"'", "-Prune")
}
