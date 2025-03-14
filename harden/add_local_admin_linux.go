package harden

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var addLocalAdminCmd = &cobra.Command{
	Use:   "add-admin [username]",
	Short: "Adds a local admin to this system",
	Long:  `Adds a local admin and prompts the user to enter the password.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		username := args[0]
		addLocalAdmin(username)
	},
}

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

func setupAddLocalAdminCmd(cmd *cobra.Command) {
	cmd.AddCommand(addLocalAdminCmd)
}
