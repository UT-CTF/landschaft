package harden

import (
	"fmt"
	"os"
	"os/exec"
)

func addLocalAdmin(username string) {
	cmd := exec.Command("useradd", username, "--no-create-home", "--home-dir", "/")
	if err := cmd.Run(); err != nil {
		fmt.Println("Error w/ useradd:", err)
		return
	}
	cmd = exec.Command("usermod", "-aG", "0", username)
	if err := cmd.Run(); err != nil {
		fmt.Println("Error w/ usermod:", err)
		return
	}
	cmd = exec.Command("passwd", username)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Error w/ passwd:", err)
	}

}
