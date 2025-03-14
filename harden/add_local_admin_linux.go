package harden

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func addLocalAdmin(username string) {
	cmd := exec.Command("useradd", username, "--no-create-home", "--home-dir", "/")
	if err := runCommand(cmd); err != nil {
		return
	}
	// Add user to group 0 (root)
	cmd = exec.Command("usermod", "-aG", "0", username)
	if err := runCommand(cmd); err != nil {
		return
	}
	// Set the user's password interactively
	cmd = exec.Command("passwd", username)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Error with passwd:", err)
	}

	fmt.Printf("Successfully added: %s\n", username)
}

func runCommand(cmd *exec.Cmd) error {
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		fmt.Printf("\nCommand failed: %s\n", cmd.String())
		fmt.Printf("Error: %v\n", err)
		fmt.Printf("Stdout: %s\n", stdout.String())
		fmt.Printf("Stderr: %s\n", stderr.String())
	}
	return err
}
