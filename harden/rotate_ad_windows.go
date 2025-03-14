package harden

import (
	"fmt"
	"strings"

	"github.com/UT-CTF/landschaft/embed"
)

func getDomainUsers() ([]string, error) {
	scriptOut, err := embed.ExecuteScript("harden/get_ad_user_list.ps1", false, "")
	if err != nil {
		return nil, fmt.Errorf("error getting user list: %w", err)
	}
	users := strings.Split(strings.TrimSpace(scriptOut), "\n")
	for i, user := range users {
		users[i] = strings.TrimSpace(user)
	}
	return users, nil
}

func applyDomainPasswordChanges(csvPath string) {
	fmt.Println("Applying password changes")
	scriptOut, err := embed.ExecuteScript("harden/rotate_ad.ps1", false, fmt.Sprintf("-Path '%s'", csvPath))
	if err != nil {
		fmt.Println("Error applying password changes:", err)
	}
	fmt.Println(scriptOut)
}
