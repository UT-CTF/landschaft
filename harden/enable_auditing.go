package harden

import (
	"fmt"

	"github.com/spf13/cobra"
)

func setupEnableAuditingCmd(cmd *cobra.Command) {
	c := &cobra.Command{
		Use:   "enable-auditing",
		Short: "Enable/verify host auditing (Windows: audit policy; Linux: auth/audit)",
		Run: func(cmd *cobra.Command, args []string) {
			if PlanMode {
				fmt.Println("Plan: would enable audit policy (Windows) or ensure auth/audit logging (Linux).")
				return
			}
			runEnableAuditing()
		},
	}
	cmd.AddCommand(c)
}
