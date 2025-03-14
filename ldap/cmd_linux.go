package ldap

import (
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

// Default character set for password generation
const defaultPasswordChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_.?!+=:^()"

// setupGeneratePasswordsCmd adds the generate-passwords subcommand
func setupGeneratePasswordsCmd(cmd *cobra.Command) {
	var (
		baseDn       string
		outputPath   string
		passwordLen  uint
		allowedChars string
		ldapArgs     []string
		excludeUsers []string
	)

	generatePasswordsCmd := &cobra.Command{
		Use:   "gen-passwords",
		Short: "Generate a CSV file with usernames and new random passwords for LDAP users",
		Long:  `Generate a CSV file with usernames and new random passwords for LDAP users. Will not overwrite existing files.`,
		Run: func(cmd *cobra.Command, args []string) {
			err := GeneratePasswordsCSV(baseDn, outputPath, passwordLen, allowedChars, ldapArgs, excludeUsers)
			if err != nil {
				log.Error("Failed to generate password:", "err", err)
				return
			}
			fmt.Printf("Successfully wrote passwords to %s\n", outputPath)
		},
	}

	// Add flags
	generatePasswordsCmd.Flags().StringVarP(&baseDn, "domain", "d", "", "Base DN (e.g. ou=Users,dc=mydom,dc=com)")
	generatePasswordsCmd.Flags().StringVar(&outputPath, "output", "new_ldap_passwords.csv", "Path to output CSV file")
	generatePasswordsCmd.Flags().UintVar(&passwordLen, "length", 16, "Length of generated passwords")
	generatePasswordsCmd.Flags().StringVar(&allowedChars, "chars", defaultPasswordChars,
		"Characters to use for password generation")
	generatePasswordsCmd.Flags().StringArrayVar(&ldapArgs, "ldap-arg", []string{},
		"Additional arguments to pass to ldapsearch (can be specified multiple times)")
	generatePasswordsCmd.Flags().StringArrayVar(&excludeUsers, "exclude", []string{},
		"Exclude specific users (can be specified multiple times). Note: 'blackteam' and 'root' are always excluded.")

	generatePasswordsCmd.MarkFlagRequired("domain")

	cmd.AddCommand(generatePasswordsCmd)
}
