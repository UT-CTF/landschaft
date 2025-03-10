package audit

import (
	"fmt"
	"os"
	"bufio"
)

func check_sshd(){
	file, err := os.Open("/etc/ssh/sshd_config")
	if err != nil {
		fmt.Println("Error reading sshd_config release:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	fmt.Println("listing sshd_config enabled rules:")
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 && line[0] != '#' {
			fmt.Println(line)
		}
	}

	scanner = bufio.NewScanner(file)
	for scanner.Scan() {
		
	}
}