package harden

import (
	"fmt"

	"github.com/UT-CTF/landschaft/util"
	"github.com/spf13/cobra"
)

// PlanMode is set when --plan is passed; subcommands then print intended actions only.
var PlanMode bool

func SetupCommand(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVar(&PlanMode, "plan", false, "Show planned changes only, do not apply")
	setupConfigureShellCmd(cmd)
	setupRotatePwdCmd(cmd)
	setupFirewallCmd(cmd)
	setupAddLocalAdminCmd(cmd)
	setupBackupEtcCmd(cmd)
	setupRestoreEtcCmd(cmd)
	setupCCDCCommands(cmd)
	setupEnableAuditingCmd(cmd)
}

func Run(cmd *cobra.Command) {
	if PlanMode {
		fmt.Println(util.TitleColor.Render("Plan mode: suggested CCDC hardening subcommands"))
		fmt.Println("  landschaft harden rotate-pwd --generate -f passwords.csv")
		fmt.Println("  landschaft harden rotate-pwd --apply -f passwords.csv")
		fmt.Println("  landschaft harden firewall --inbound --apply -f rules.json")
		fmt.Println("  landschaft harden backup-etc")
		fmt.Println("  landschaft harden add-local-admin <user>")
		fmt.Println("  landschaft harden configure-shell (Linux: shell logging)")
		fmt.Println()
	}
	fmt.Println(util.ErrorStyle.Render("Error: No subcommand specified"))
	fmt.Println()
	_ = cmd.Usage()
}
