package ldap

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	ldifTemplatePath string
	csvInputPath     string
	outputPath       string
)

var generateLdifCommand = &cobra.Command{
	Use:   "ldif",
	Short: "Generate LDIF files",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Check if required flags are provided
		if ldifTemplatePath == "" {
			return fmt.Errorf("template flag is required")
		}
		if csvInputPath == "" {
			return fmt.Errorf("input flag is required")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := generateLdif(ldifTemplatePath, csvInputPath, outputPath); err != nil {
			fmt.Println(err)
		}
	},
}

func setupGenerateLdifCmd(cmd *cobra.Command) {
	// Add flags
	generateLdifCommand.Flags().StringVarP(&ldifTemplatePath, "template", "t", "", "Path to LDIF template file (required)")
	generateLdifCommand.Flags().StringVarP(&csvInputPath, "input", "i", "", "Path to input CSV file (required)")
	generateLdifCommand.Flags().StringVarP(&outputPath, "output", "o", "output.ldif", "Path to output LDIF file")

	// Mark flags as required
	generateLdifCommand.MarkFlagRequired("template")
	generateLdifCommand.MarkFlagRequired("input")

	cmd.AddCommand(generateLdifCommand)
}
