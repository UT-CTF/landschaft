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

var (
	apply, generate   bool
	length            uint
	blacklist         []string
	allowedCharacters string
	filePath          string
	domain            bool
)

var rotatePwdCmd = &cobra.Command{
	Use:   "rotate-pwd",
	Short: "Rotate user passwords",
	Long: `Rotate users passwords in two steps:
"Generate" will generate a csv of all new passwords.
"Apply" will set all passwords to the new passwords.`,
	Run: func(cmd *cobra.Command, args []string) {
		getUsersCmd := getLocalUsers
		applyPasswordCmd := applyPasswordChanges
		if domain {
			getUsersCmd = getDomainUsers
			applyPasswordCmd = applyDomainPasswordChanges
		}
		rotateLocalUsers(apply, generate, filePath, length, blacklist, allowedCharacters, getUsersCmd, applyPasswordCmd)
	},
}

func setupRotatePwdCmd(cmd *cobra.Command) {
	rotatePwdCmd.Flags().BoolVar(&apply, "apply", false, "Apply the new passwords to the users")
	rotatePwdCmd.Flags().BoolVarP(&generate, "generate", "g", false, "Generate new passwords for the users")
	rotatePwdCmd.Flags().UintVar(&length, "length", 16, "Length of the new passwords")
	rotatePwdCmd.Flags().StringSliceVar(&blacklist, "blacklist", []string{}, fmt.Sprintf("Additional list of users to exclude from password rotation (default: %s)", strings.Join(defaultBlacklist, ", ")))
	rotatePwdCmd.Flags().StringVarP(&filePath, "file", "f", "", "File path to save the new passwords or read generated passwords (csv format)")
	rotatePwdCmd.Flags().StringVar(&allowedCharacters, "allowed-chars", defaultAllowedCharacters, "Allowed characters for the new passwords")
	rotatePwdCmd.Flags().BoolVar(&domain, "domain", false, "Rotate domain users instead of local users (AD on Windows, LDAP on Linux)")
	rotatePwdCmd.MarkFlagRequired("file")
	rotatePwdCmd.MarkFlagsMutuallyExclusive("apply", "generate")
	rotatePwdCmd.MarkFlagsOneRequired("apply", "generate")

	cmd.AddCommand(rotatePwdCmd)
}

func rotateLocalUsers(apply bool, generate bool, filePath string, length uint, blacklist []string, allowedCharacters string, getUsersCmd func() ([]string, error), applyPasswordCmd func(string)) {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		fmt.Println("Error getting absolute path")
		return
	}
	if generate {
		users, err := getUsersCmd()
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
		applyPasswordCmd(absPath)
	}
}

var defaultAllowedCharacters string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_.?!+=:^()"

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
