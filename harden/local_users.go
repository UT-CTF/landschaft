package harden

import (
	"fmt"
	"os"

	"github.com/UT-CTF/landschaft/util"
)

func generatePasswordChangeCSV(users []string, length uint, filePath string, strict bool) {
	passwords := make([]string, len(users))
	for i := range passwords {
		if strict {
			passwords[i] = util.GenerateStrictRandomPassword(length)
		} else {
			passwords[i] = util.GenerateRandomPassword(length)
		}
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
