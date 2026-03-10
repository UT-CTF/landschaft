/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"time"

	"github.com/UT-CTF/landschaft/audit"
	"github.com/spf13/cobra"
)

var auditArgs struct {
	sshd     bool
	versions bool
	max      int
	timeout  time.Duration
}

// auditCmd represents the audit command
var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Audit a host for common security issues",
	Run: func(cmd *cobra.Command, args []string) {
		audit.Run(audit.Options{
			CheckSSHD:     auditArgs.sshd,
			CheckVersions: auditArgs.versions,
			MaxPackages:   auditArgs.max,
			Timeout:       auditArgs.timeout,
		})
	},
}

func init() {
	rootCmd.AddCommand(auditCmd)

	def := audit.DefaultOptions()
	auditCmd.Flags().BoolVar(&auditArgs.sshd, "sshd", def.CheckSSHD, "Audit SSHD configuration (Linux only)")
	auditCmd.Flags().BoolVar(&auditArgs.versions, "versions", def.CheckVersions, "Check installed software versions against OSV (best-effort)")
	auditCmd.Flags().IntVar(&auditArgs.max, "max", def.MaxPackages, "Maximum packages to query against OSV")
	auditCmd.Flags().DurationVar(&auditArgs.timeout, "timeout", def.Timeout, "Timeout for OSV queries (e.g. 5s)")

}
