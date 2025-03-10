package harden

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/UT-CTF/landschaft/util"
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
			return
		}
		generate, err := cmd.Flags().GetBool("generate")
		if err != nil {
			fmt.Println("Error getting generate flag")
			return
		}
		filePath, err := cmd.Flags().GetString("file")
		if err != nil {
			fmt.Println("Error getting file path")
			return
		}
		length, err := cmd.Flags().GetUint("length")
		if err != nil {
			fmt.Println("Error getting length flag")
			return
		}
		blacklist, err := cmd.Flags().GetStringSlice("blacklist")
		if err != nil {
			fmt.Println("Error getting blacklist flag")
			return
		}
		allowedCharacters, err := cmd.Flags().GetString("allowed-chars")
		if err != nil {
			fmt.Println("Error getting allowed-chars flag")
			return
		}
		rotateLocalUsers(apply, generate, filePath, length, blacklist, allowedCharacters)
	},
}

func rotateLocalUsers(apply bool, generate bool, filePath string, length uint, blacklist []string, allowedCharacters string) {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		fmt.Println("Error getting absolute path")
		return
	}
	if generate {
		users, err := getLocalUsers()
		if err != nil {
			fmt.Println(err)
			return
		}
		filtered, removed := filterBlacklistedUsers(users, append(defaultBlacklist, blacklist...))
		if len(removed) > 0 {
			fmt.Printf("Blacklisted users removed:\n%s\n", strings.Join(removed, "\n"))
		}
		fmt.Printf("Changing Passwords for %d users:\n%s\n", len(filtered), strings.Join(filtered, "\n"))
		generatePasswordChangeCSV(filtered, length, absPath, true, allowedCharacters)
	} else if apply {
		applyPasswordChanges(absPath)
	}
}

var defaultAllowedCharacters string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_.?!+=:^()"

func setupRotateLocalUsersCmd(cmd *cobra.Command) {
	rotateLocalUsersCmd.Flags().Bool("apply", false, "Apply the new passwords to the users")
	rotateLocalUsersCmd.Flags().BoolP("generate", "g", false, "Generate new passwords for the users")
	rotateLocalUsersCmd.Flags().Uint("length", 16, "Length of the new passwords")
	rotateLocalUsersCmd.Flags().StringSlice("blacklist", []string{}, fmt.Sprintf("Additional list of users to exclude from password rotation (default: %s)", strings.Join(defaultBlacklist, ", ")))
	rotateLocalUsersCmd.Flags().StringP("file", "f", "", "File path to save the new passwords or read generated passwords (csv format)")
	rotateLocalUsersCmd.Flags().String("allowed-chars", defaultAllowedCharacters, "Allowed characters for the new passwords")
	rotateLocalUsersCmd.MarkFlagRequired("file")
	rotateLocalUsersCmd.MarkFlagsMutuallyExclusive("apply", "generate")
	rotateLocalUsersCmd.MarkFlagsOneRequired("apply", "generate")

	cmd.AddCommand(rotateLocalUsersCmd)
}

func filterBlacklistedUsers(users []string, blacklist []string) ([]string, []string) {
	filteredUsers := make([]string, 0)
	removedUsers := make([]string, 0)
	for _, user := range users {
		if slices.Contains(blacklist, user) {
			removedUsers = append(removedUsers, user)
		} else {
			filteredUsers = append(filteredUsers, user)
		}
	}
	return filteredUsers, removedUsers
}

func generatePasswordChangeCSV(users []string, length uint, filePath string, strict bool, allowedCharacters string) {
	passwords := make([]string, len(users))
	for i := range passwords {
		passwords[i] = util.GenerateRandomPassword(length, allowedCharacters, strict)
	}

	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file")
		return
	}
	defer file.Close()

	for i, user := range users {
		_, err := file.WriteString(fmt.Sprintf("%s,%s\n", user, passwords[i]))
		if err != nil {
			fmt.Println("Error writing to file")
			return
		}
	}

	fmt.Println("Wrote passwords csv to", filePath)
}
