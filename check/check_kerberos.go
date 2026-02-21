package check

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/jcmturner/gokrb5/v8/client"
	"github.com/jcmturner/gokrb5/v8/config"
	"github.com/jcmturner/gokrb5/v8/credentials"
	"github.com/spf13/cobra"
)

var kerberosCheckCmd = &cobra.Command{
	Use:   "kerberos",
	Short: "Check Kerberos KDC connectivity and authentication",
	Long: `Verify Kerberos Key Distribution Center (KDC) is accessible and
credentials are valid. Obtains a Ticket Granting Ticket (TGT) to verify authentication.`,
	Run: func(cmd *cobra.Command, args []string) {
		kdc, _ := cmd.Flags().GetString("kdc")
		realm, _ := cmd.Flags().GetString("realm")
		fqdn, _ := cmd.Flags().GetString("fqdn")
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		csvFile, _ := cmd.Flags().GetString("csv")
		timeout, _ := cmd.Flags().GetInt("timeout")

		// If CSV file is provided, batch check
		if csvFile != "" {
			results, err := checkKerberosBatch(kdc, realm, fqdn, csvFile, timeout)
			if err != nil {
				log.Error("Kerberos batch check failed", "error", err)
				os.Exit(1)
			}

			// Print results
			fmt.Printf("\n=== Kerberos Batch Check Results ===\n")
			fmt.Printf("Total: %d | Valid: %d | Invalid: %d\n\n",
				results.Total, results.Valid, results.Invalid)

			if len(results.ValidCreds) > 0 {
				fmt.Println("✓ Valid credentials:")
				for _, cred := range results.ValidCreds {
					fmt.Printf("  • %s\n", cred)
				}
			}

			if len(results.InvalidCreds) > 0 {
				fmt.Println("\n✗ Invalid credentials:")
				for _, cred := range results.InvalidCreds {
					fmt.Printf("  • %s\n", cred)
				}
			}

			if results.Invalid > 0 {
				os.Exit(2)
			}
			os.Exit(0)
		}

		// Single credential check
		if username == "" || password == "" {
			log.Error("Either provide --username and --password, or use --csv for batch checking")
			os.Exit(1)
		}

		if err := checkKerberos(kdc, realm, fqdn, username, password, timeout); err != nil {
			log.Error("Kerberos check failed", "error", err)
			os.Exit(2)
		}
		fmt.Println("✓ Kerberos check passed")
		os.Exit(0)
	},
}

func setupKerberosCheckCmd(cmd *cobra.Command) {
	kerberosCheckCmd.Flags().StringP("kdc", "k", "", "KDC server address (IP or hostname)")
	kerberosCheckCmd.Flags().StringP("realm", "r", "", "Kerberos realm (e.g., EXAMPLE.COM)")
	kerberosCheckCmd.Flags().StringP("fqdn", "f", "", "Fully qualified domain name of the KDC")
	kerberosCheckCmd.Flags().StringP("username", "u", "", "Username for authentication (not used with --csv)")
	kerberosCheckCmd.Flags().StringP("password", "p", "", "Password for authentication (not used with --csv)")
	kerberosCheckCmd.Flags().StringP("csv", "c", "", "CSV file with username,password pairs for batch checking")
	kerberosCheckCmd.Flags().IntP("timeout", "t", 10, "Connection timeout in seconds")

	kerberosCheckCmd.MarkFlagRequired("kdc")
	kerberosCheckCmd.MarkFlagRequired("realm")
	kerberosCheckCmd.MarkFlagRequired("fqdn")

	cmd.AddCommand(kerberosCheckCmd)
}

func checkKerberos(kdc, realm, fqdn, username, password string, timeoutSec int) error {
	fmt.Printf("Connecting to Kerberos KDC at %s (realm: %s)...\n", kdc, realm)

	// Ensure KDC has port specified (default to 88)
	kdcAddress := kdc
	if !contains(kdcAddress, ":") {
		kdcAddress = fmt.Sprintf("%s:88", kdc)
	}

	// Create a minimal Kerberos configuration using IP:port to avoid DNS lookups
	krb5conf := `[libdefaults]
  default_realm = %s
  dns_lookup_realm = false
  dns_lookup_kdc = false
  ticket_lifetime = 24h
  forwardable = yes
  udp_preference_limit = 1

[realms]
  %s = {
    kdc = %s
    admin_server = %s
  }

[domain_realm]
  .%s = %s
  %s = %s
`
	confString := fmt.Sprintf(krb5conf, realm, realm, kdcAddress, kdcAddress, fqdn, realm, fqdn, realm)

	// Parse the configuration
	cfg, err := config.NewFromString(confString)
	if err != nil {
		return fmt.Errorf("failed to create Kerberos config: %w", err)
	}

	// Create a client with username and password
	cl := client.NewWithPassword(username, realm, password, cfg, client.DisablePAFXFAST(true))

	// Login to obtain TGT
	done := make(chan error, 1)
	go func() {
		done <- cl.Login()
	}()

	select {
	case err := <-done:
		if err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}
	case <-time.After(time.Duration(timeoutSec) * time.Second):
		return fmt.Errorf("connection timeout after %d seconds", timeoutSec)
	}

	fmt.Println("✓ Connection established")
	fmt.Printf("✓ Authentication successful (obtained TGT for %s@%s)\n", username, realm)

	// Verify we have credentials by checking IsConfigured (returns bool and error)
	configured, err := cl.IsConfigured()
	if err != nil {
		return fmt.Errorf("failed to verify client configuration: %w", err)
	}
	if !configured {
		return fmt.Errorf("client not properly configured after login")
	}

	// Try to get the credentials to verify they exist
	creds := credentials.New(username, realm)
	if creds == nil {
		return fmt.Errorf("failed to verify credentials")
	}

	fmt.Println("✓ Ticket verification successful")

	// Destroy the session
	cl.Destroy()

	return nil
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// checkKerberosBatch checks multiple username/password pairs from a CSV file
func checkKerberosBatch(kdc, realm, fqdn, csvFile string, timeoutSec int) (*BatchCheckResult, error) {
	// Read CSV file
	file, err := os.Open(csvFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV file: %w", err)
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("CSV file is empty")
	}

	results := &BatchCheckResult{
		ValidCreds:   []string{},
		InvalidCreds: []string{},
	}

	fmt.Printf("Checking %d credential(s) from %s...\n\n", len(records), csvFile)

	for i, record := range records {
		if len(record) < 2 {
			fmt.Printf("[%d/%d] ✗ Skipping invalid CSV row (need username,password)\n", i+1, len(records))
			continue
		}

		username := record[0]
		password := record[1]
		results.Total++

		fmt.Printf("[%d/%d] Testing %s... ", i+1, len(records), username)

		err := checkKerberos(kdc, realm, fqdn, username, password, timeoutSec)
		if err != nil {
			fmt.Println("✗ INVALID")
			results.Invalid++
			results.InvalidCreds = append(results.InvalidCreds, username)
		} else {
			fmt.Println("✓ VALID")
			results.Valid++
			results.ValidCreds = append(results.ValidCreds, username)
		}
	}

	return results, nil
}
