package harden

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func addLocalAdmin(username string) {
	// Create the user
	cmd := exec.Command("sudo", "useradd", username, "--no-create-home", "--home-dir", "/")
	if err := runCommand(cmd); err != nil {
		return
	}

	// Add user to group 0 (root)
	cmd = exec.Command("sudo", "usermod", "-aG", "0", username)
	if err := runCommand(cmd); err != nil {
		return
	}

	// Set the user's password interactively
	cmd = exec.Command("sudo", "passwd", username)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Error with passwd:", err)
	}

}

func runCommand(cmd *exec.Cmd) error {
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		fmt.Printf("\n❌ Command failed: %s\n", cmd.String())
		fmt.Printf("Error: %v\n", err)
		fmt.Printf("Stdout: %s\n", stdout.String())
		fmt.Printf("Stderr: %s\n", stderr.String())
	}
	return err
}
