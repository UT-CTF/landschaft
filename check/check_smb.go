package check

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/hirochachacha/go-smb2"
	"github.com/spf13/cobra"
)

var smbCheckCmd = &cobra.Command{
	Use:   "smb",
	Short: "Check SMB server connectivity and authentication",
	Long: `Verify SMB server is accessible and credentials are valid.
Can list shares or access a specific file on a share.`,
	Run: func(cmd *cobra.Command, args []string) {
		host, _ := cmd.Flags().GetString("host")
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		domain, _ := cmd.Flags().GetString("domain")
		share, _ := cmd.Flags().GetString("share")
		path, _ := cmd.Flags().GetString("path")
		timeout, _ := cmd.Flags().GetInt("timeout")

		if err := checkSMB(host, username, password, domain, share, path, timeout); err != nil {
			log.Error("SMB check failed", "error", err)
			os.Exit(2)
		}
		fmt.Println("✓ SMB check passed")
		os.Exit(0)
	},
}

func setupSmbCheckCmd(cmd *cobra.Command) {
	smbCheckCmd.Flags().StringP("host", "H", "", "SMB server address (IP or hostname)")
	smbCheckCmd.Flags().StringP("username", "u", "", "Username for authentication")
	smbCheckCmd.Flags().StringP("password", "p", "", "Password for authentication")
	smbCheckCmd.Flags().StringP("domain", "d", "", "Domain name (optional, defaults to WORKGROUP)")
	smbCheckCmd.Flags().StringP("share", "s", "", "Share name to list (e.g., C$, IPC$)")
	smbCheckCmd.Flags().String("path", "", "Path to a file on the share to access (e.g., /Windows/System32/config)")
	smbCheckCmd.Flags().IntP("timeout", "t", 10, "Connection timeout in seconds")

	smbCheckCmd.MarkFlagRequired("host")
	smbCheckCmd.MarkFlagRequired("username")
	smbCheckCmd.MarkFlagRequired("password")

	cmd.AddCommand(smbCheckCmd)
}

func checkSMB(host, username, password, domain, share, path string, timeoutSec int) error {
	// Default domain if not specified
	if domain == "" {
		domain = "WORKGROUP"
	}

	// Connect to SMB server on port 445
	address := fmt.Sprintf("%s:445", host)
	fmt.Printf("Connecting to SMB server at %s...\n", address)

	// Set timeout for dial
	conn, err := net.DialTimeout("tcp", address, time.Duration(timeoutSec)*time.Second)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer conn.Close()

	fmt.Println("✓ Connection established")

	// Create SMB session
	d := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     username,
			Password: password,
			Domain:   domain,
		},
	}

	session, err := d.Dial(conn)
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}
	defer session.Logoff()

	fmt.Printf("✓ Authentication successful (domain: %s, user: %s)\n", domain, username)

	// If share is specified, try to mount it
	if share != "" {
		fmt.Printf("Mounting share: %s\n", share)
		fs, err := session.Mount(share)
		if err != nil {
			return fmt.Errorf("failed to mount share: %w", err)
		}
		defer fs.Umount()

		fmt.Printf("✓ Share mounted successfully: %s\n", share)

		// If path is specified, try to access it
		if path != "" {
			fmt.Printf("Accessing path: %s\n", path)
			stat, err := fs.Stat(path)
			if err != nil {
				return fmt.Errorf("failed to access path: %w", err)
			}

			if stat.IsDir() {
				fmt.Printf("✓ Path accessible (directory): %s\n", path)
			} else {
				fmt.Printf("✓ Path accessible (file, size: %d bytes): %s\n", stat.Size(), path)
			}
		} else {
			// List the share contents (root level)
			entries, err := fs.ReadDir(".")
			if err != nil {
				return fmt.Errorf("failed to list share contents: %w", err)
			}

			fmt.Printf("✓ Share listing successful (%d entries in root)\n", len(entries))
			if len(entries) > 0 {
				fmt.Println("\nFirst few entries:")
				for i, entry := range entries {
					if i >= 5 {
						break
					}
					entryType := "file"
					if entry.IsDir() {
						entryType = "dir "
					}
					fmt.Printf("  [%s] %s\n", entryType, entry.Name())
				}
			}
		}
	} else {
		fmt.Println("✓ Session verification successful (no share specified)")
	}

	return nil
}
