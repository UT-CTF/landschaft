package cmd

import (
	"github.com/UT-CTF/landschaft/check"
	"github.com/spf13/cobra"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Perform service checks to verify connectivity and authentication",
	Long: `Check various services (LDAP, Kerberos, SMB) to verify they are up and
functioning as intended. Each subcommand tests authentication and optionally
performs additional operations to validate service health.`,
	Run: func(cmd *cobra.Command, args []string) {
		check.Run(cmd)
	},
}

func init() {
	check.SetupCommand(checkCmd)
	rootCmd.AddCommand(checkCmd)
}
