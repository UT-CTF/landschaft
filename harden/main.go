package harden

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rotateLocalUsersCmd = &cobra.Command{
	Use:   "rotate-local [file path]",
	Short: "Rotate local user passwords",
	Long: `Rotate local users passwords in two steps:
"Generate" will generate a csv of all new passwords. 
"Apply" will set all passwords to the new passwords.`,
	Run: func(cmd *cobra.Command, args []string) {
		apply, err := cmd.Flags().GetBool("apply")
		if err != nil {
			fmt.Println("Error getting apply flag")
		}
		generate, err := cmd.Flags().GetBool("generate")
		if err != nil {
			fmt.Println("Error getting generate flag")
		}
		filePath, err := cmd.Flags().GetString("file")
		if err != nil {
			fmt.Println("Error getting file path")
		}
		length, err := cmd.Flags().GetUint("length")
		if err != nil {
			fmt.Println("Error getting length flag")
		}
		rotateLocalUsers(apply, generate, filePath, length)
	},
}

func SetupCommand(cmd *cobra.Command) {
	rotateLocalUsersCmd.Flags().Bool("apply", false, "Apply the new passwords to the users")
	rotateLocalUsersCmd.Flags().BoolP("generate", "g", false, "Generate new passwords for the users")
	rotateLocalUsersCmd.Flags().Uint("length", 16, "Length of the new passwords")
	rotateLocalUsersCmd.Flags().StringP("file", "f", "", "File path to save the new passwords or read generated passwords (csv format)")
	rotateLocalUsersCmd.MarkFlagRequired("file")
	rotateLocalUsersCmd.MarkFlagsMutuallyExclusive("apply", "generate")
	rotateLocalUsersCmd.MarkFlagsOneRequired("apply", "generate")

	cmd.AddCommand(rotateLocalUsersCmd)
}

func Run(cmd *cobra.Command) {
	fmt.Println("Error: No subcommand specified")
	fmt.Println()
	_ = cmd.Usage()
}
