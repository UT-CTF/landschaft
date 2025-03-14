package ldap

import (
	"github.com/spf13/cobra"
)

// setupGeneratePasswordsCmd is a no-op on Windows
func setupGeneratePasswordsCmd(cmd *cobra.Command) {
	// Command not available on Windows
}
