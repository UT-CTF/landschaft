package ldap

import (
	"fmt"

	"github.com/spf13/cobra"
)

var generateLdifSingleCommand = &cobra.Command{
	Use:   "ldifsingle",
	Short: "Generate LDIF files for a single user/password pair",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Check if required flags are provided
		if ldifTemplatePath == "" {
			return fmt.Errorf("template flag is required")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := generateLdifSingle(ldifTemplatePath, outputPath); err != nil {
			fmt.Println(err)
		}
	},
}

func setupGenerateLdifSingleCmd(cmd *cobra.Command) {
	// Add flags
	generateLdifSingleCommand.Flags().StringVarP(&ldifTemplatePath, "template", "t", "", "Path to LDIF template file (required)")
	generateLdifSingleCommand.Flags().StringVarP(&outputPath, "output", "o", "output.ldif", "Path to output LDIF file")

	// Mark flags as required
	generateLdifSingleCommand.MarkFlagRequired("template")

	cmd.AddCommand(generateLdifSingleCommand)
}
