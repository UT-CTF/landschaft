package harden

import (
	"fmt"
	"strings"

	"github.com/UT-CTF/landschaft/embed"
)

func rotateLocalUsers(apply bool, generate bool, filePath string, length uint) {
	if generate {
		scriptOut, err := embed.ExecuteScript("harden/get_user_list.ps1", "")
		if err != nil {
			fmt.Println("Error getting user list")
			return
		}
		users := strings.Split(strings.TrimSpace(scriptOut), "\n")
		for i, user := range users {
			users[i] = strings.TrimSpace(user)
		}
		generatePasswordChangeCSV(users, length, filePath, true)
	}
}
