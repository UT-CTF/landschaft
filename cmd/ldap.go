package cmd

import (
	"github.com/spf13/cobra"

	"github.com/UT-CTF/landschaft/ldap"
)

var ldapCmd = &cobra.Command{
	Use:   "ldap",
	Short: "LDIF templating",
	Run: func(cmd *cobra.Command, args []string) {
		ldap.Run(cmd)
	},
}

func init() {
	ldap.SetupCommand(ldapCmd)

	rootCmd.AddCommand(ldapCmd)
}
