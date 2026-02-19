package check

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

var sshCheckCmd = &cobra.Command{
	Use:   "ssh",
	Short: "Check SSH server connectivity and authentication",
	Long: `Verify SSH server is accessible and credentials are valid.
Optionally executes a command to verify full session functionality.`,
	Run: func(cmd *cobra.Command, args []string) {
		host, _ := cmd.Flags().GetString("host")
		port, _ := cmd.Flags().GetInt("port")
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		keyFile, _ := cmd.Flags().GetString("key")
		command, _ := cmd.Flags().GetString("command")
		timeout, _ := cmd.Flags().GetInt("timeout")

		if err := checkSSH(host, port, username, password, keyFile, command, timeout); err != nil {
			log.Error("SSH check failed", "error", err)
			os.Exit(2)
		}
		fmt.Println("✓ SSH check passed")
		os.Exit(0)
	},
}

func setupSshCheckCmd(cmd *cobra.Command) {
	sshCheckCmd.Flags().StringP("host", "H", "", "SSH server address (IP or hostname)")
	sshCheckCmd.Flags().IntP("port", "P", 22, "SSH server port")
	sshCheckCmd.Flags().StringP("username", "u", "", "Username for authentication")
	sshCheckCmd.Flags().StringP("password", "p", "", "Password for authentication")
	sshCheckCmd.Flags().StringP("key", "k", "", "Path to private key file for authentication")
	sshCheckCmd.Flags().StringP("command", "c", "", "Optional command to execute (e.g., 'whoami')")
	sshCheckCmd.Flags().IntP("timeout", "t", 10, "Connection timeout in seconds")

	sshCheckCmd.MarkFlagRequired("host")
	sshCheckCmd.MarkFlagRequired("username")

	cmd.AddCommand(sshCheckCmd)
}

func checkSSH(host string, port int, username, password, keyFile, command string, timeoutSec int) error {
	address := fmt.Sprintf("%s:%d", host, port)
	fmt.Printf("Connecting to SSH server at %s...\n", address)

	// Prepare authentication methods
	var authMethods []ssh.AuthMethod

	// Add password authentication if provided
	if password != "" {
		authMethods = append(authMethods, ssh.Password(password))
	}

	// Add public key authentication if key file is provided
	if keyFile != "" {
		key, err := os.ReadFile(keyFile)
		if err != nil {
			return fmt.Errorf("failed to read private key file: %w", err)
		}

		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return fmt.Errorf("failed to parse private key: %w", err)
		}

		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}

	if len(authMethods) == 0 {
		return fmt.Errorf("no authentication method provided (need --password or --key)")
	}

	// Configure SSH client
	config := &ssh.ClientConfig{
		User:            username,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // For checking, we skip host key verification
		Timeout:         time.Duration(timeoutSec) * time.Second,
	}

	// Connect to SSH server
	client, err := ssh.Dial("tcp", address, config)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer client.Close()

	fmt.Println("✓ Connection established")
	fmt.Printf("✓ Authentication successful (user: %s)\n", username)

	// If a command is specified, execute it
	if command != "" {
		session, err := client.NewSession()
		if err != nil {
			return fmt.Errorf("failed to create session: %w", err)
		}
		defer session.Close()

		fmt.Printf("Executing command: %s\n", command)
		output, err := session.CombinedOutput(command)
		if err != nil {
			return fmt.Errorf("command execution failed: %w", err)
		}

		fmt.Printf("✓ Command executed successfully\nOutput:\n%s\n", string(output))
	} else {
		// Just verify session creation works
		session, err := client.NewSession()
		if err != nil {
			return fmt.Errorf("failed to create session: %w", err)
		}
		session.Close()
		fmt.Println("✓ Session verification successful")
	}

	return nil
}
