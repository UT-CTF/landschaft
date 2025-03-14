package harden

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/UT-CTF/landschaft/embed"
	"github.com/UT-CTF/landschaft/util"
	"github.com/rivo/tview"
)

type firewallRule struct {
	Name      string `json:"Name"`
	Direction string `json:"Direction"`
	Action    string `json:"Action"`
	Protocol  string `json:"Protocol"`
	LocalPort string `json:"LocalPort"`
	Program   string `json:"Program"`
}

func configureFirewall() {
	jsonStr, err := embed.ExecuteScript("harden/firewall.ps1", false)
	if err != nil {
		fmt.Println("Error executing script: ", err)
		return
	}
	var rules []firewallRule
	err = json.Unmarshal([]byte(jsonStr), &rules)
	if err != nil {
		fmt.Println("Error parsing JSON: ", err)
		return
	}

	// for _, rule := range rules {
	// 	fmt.Printf("Adding rule %s\n", rule.Name)
	// 	fmt.Printf("Direction: %s\n", rule.Direction)
	// 	fmt.Printf("Action: %s\n", rule.Action)
	// 	fmt.Printf("Protocol: %s\n", rule.Protocol)
	// 	fmt.Printf("LocalPort: %s\n", rule.LocalPort)
	// 	fmt.Printf("Program: %s\n", rule.Program)
	// }

	app := tview.NewApplication()
	list := tview.NewList()

	selectedRules := make(map[string]bool)

	runRules := false

	for i, rule := range rules {
		ruleName := rule.Name
		list.AddItem(fmt.Sprintf("[ ] %s", ruleName), "", 0, func() {
			selectedRules[ruleName] = !selectedRules[ruleName]
			if selectedRules[ruleName] {
				list.SetItemText(i, fmt.Sprintf(tview.Escape("[X] %s"), ruleName), "")
			} else {
				list.SetItemText(i, fmt.Sprintf("[ ] %s", ruleName), "")
			}
		})
	}

	list.SetTitle(" Select Firewall Rules to Enable ").SetBorder(true)

	// Submit button
	list.AddItem("Enable Selected Rules", "Press Enter to enable", 'e', func() {
		runRules = true
		app.Stop()
	})

	// Quit option
	list.AddItem("Quit", "Press Q to exit", 'q', func() {
		app.Stop()
	})

	// Run the TUI
	if err := app.SetRoot(list, true).Run(); err != nil {
		fmt.Println("Error running TUI: ", err)
		return
	}

	if !runRules {
		return
	}

	var rulesToEnable []firewallRule
	for _, rule := range rules {
		if selectedRules[rule.Name] {
			rulesToEnable = append(rulesToEnable, rule)
		}
	}

	// write to new json file
	// rulesToEnableJSON, err := json.Marshal(rulesToEnable)
	rulesToEnableJSON, err := json.MarshalIndent(rulesToEnable, "", "\t")
	if err != nil {
		fmt.Println("Error marshalling JSON: ", err)
		return
	}

	filePath := "./enable_firewall_rules.json"
	filePath, err = filepath.Abs(filePath)
	if err != nil {
		fmt.Println("Error getting absolute path: ", err)
		return
	}

	err = os.WriteFile(filePath, rulesToEnableJSON, 0644)
	defer os.Remove(filePath)
	if err != nil {
		fmt.Println("Error writing JSON to file: ", err)
		return
	}

	backupPath := "./backup_firewall_rules.wfw"
	backupPath, err = filepath.Abs(backupPath)
	if err != nil {
		fmt.Println("Error getting absolute path: ", err)
		return
	}

	fmt.Println("----------------------------------------\n\n\n\n\n\n----------------------------------------")

	fmt.Println("This will remove ALL existing firewall rules and apply the selected rules.")
	fmt.Println("This will also create a backup of the current rules.")
	fmt.Println("The new rules are in the file: ", filePath)
	fmt.Println("Are you sure you want to continue? (y/n)")
	var confirm string
	fmt.Scanln(&confirm)
	if confirm != "y" {
		fmt.Println("Aborting.")
		return
	}

	util.RunAndPrintScript("harden/firewall.ps1", "-RulePath", "'"+filePath+"'", "-BackupPath", "'"+backupPath+"'", "-Apply")
}
