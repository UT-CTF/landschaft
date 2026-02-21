package check

import (
	"crypto/tls"
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/go-ldap/ldap/v3"
	"github.com/spf13/cobra"
)

var ldapCheckCmd = &cobra.Command{
	Use:   "ldap",
	Short: "Check LDAP server connectivity and authentication",
	Long: `Verify LDAP server is accessible and credentials are valid.
Performs a simple bind operation and optionally searches the directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		server, _ := cmd.Flags().GetString("server")
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		csvFile, _ := cmd.Flags().GetString("csv")
		useTLS, _ := cmd.Flags().GetBool("tls")
		timeout, _ := cmd.Flags().GetInt("timeout")

		// If CSV file is provided, batch check
		if csvFile != "" {
			results, err := checkLDAPBatch(server, csvFile, useTLS, timeout)
			if err != nil {
				log.Error("LDAP batch check failed", "error", err)
				os.Exit(1)
			}

			// Print results
			fmt.Printf("\n=== LDAP Batch Check Results ===\n")
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

		if err := checkLDAP(server, username, password, useTLS, timeout); err != nil {
			log.Error("LDAP check failed", "error", err)
			os.Exit(2)
		}
		fmt.Println("✓ LDAP check passed")
		os.Exit(0)
	},
}

func setupLdapCheckCmd(cmd *cobra.Command) {
	ldapCheckCmd.Flags().StringP("server", "s", "", "LDAP server address (IP or hostname)")
	ldapCheckCmd.Flags().StringP("username", "u", "", "Username for authentication (not used with --csv)")
	ldapCheckCmd.Flags().StringP("password", "p", "", "Password for authentication (not used with --csv)")
	ldapCheckCmd.Flags().StringP("csv", "f", "", "CSV file with username,password pairs for batch checking")
	ldapCheckCmd.Flags().Bool("tls", false, "Use StartTLS to upgrade connection to TLS")
	ldapCheckCmd.Flags().IntP("timeout", "t", 10, "Connection timeout in seconds")

	ldapCheckCmd.MarkFlagRequired("server")

	cmd.AddCommand(ldapCheckCmd)
}

func checkLDAP(server, username, password string, useTLS bool, timeoutSec int) error {
	// Always connect on port 389 initially
	port := "389"
	address := fmt.Sprintf("%s:%s", server, port)

	if useTLS {
		fmt.Printf("Connecting to ldap://%s (will upgrade to TLS via StartTLS)...\n", address)
	} else {
		fmt.Printf("Connecting to ldap://%s...\n", address)
	}

	// Set dial timeout
	ldap.DefaultTimeout = time.Duration(timeoutSec) * time.Second

	// Connect to LDAP server
	conn, err := ldap.Dial("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer conn.Close()

	fmt.Println("✓ Connection established")

	// If TLS is requested, upgrade the connection using StartTLS
	if useTLS {
		// Configure TLS to skip certificate verification for service checks
		// In production environments, proper certificate validation should be used
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
		}

		err = conn.StartTLS(tlsConfig)
		if err != nil {
			return fmt.Errorf("failed to start TLS: %w", err)
		}
		fmt.Println("✓ TLS negotiation successful")
	}

	// First, try to discover the domain name from RootDSE
	var domain string
	searchRequest := ldap.NewSearchRequest(
		"",
		ldap.ScopeBaseObject,
		ldap.NeverDerefAliases,
		0,
		timeoutSec,
		false,
		"(objectClass=*)",
		[]string{"defaultNamingContext", "rootDomainNamingContext"},
		nil,
	)

	sr, err := conn.Search(searchRequest)
	if err == nil && len(sr.Entries) > 0 {
		entry := sr.Entries[0]
		if defaultNC := entry.GetAttributeValue("defaultNamingContext"); defaultNC != "" {
			// Extract domain from DC components (e.g., DC=office,DC=local -> office.local)
			domain = ldapDNToDomain(defaultNC)
		}
	}

	// Attempt to bind (authenticate)
	// Try various formats for Active Directory and standard LDAP
	bindDNs := []string{
		username, // Simple username
	}

	// If we discovered a domain, try Active Directory formats
	if domain != "" {
		bindDNs = append([]string{
			fmt.Sprintf("%s@%s", username, domain),                              // UserPrincipalName format (user@domain.com)
			fmt.Sprintf("%s\\%s", domain[:len(domain)-len(".local")], username), // DOMAIN\user format (for .local domains)
		}, bindDNs...)
	}

	// Add traditional LDAP DN formats
	bindDNs = append(bindDNs,
		fmt.Sprintf("cn=%s", username),
		fmt.Sprintf("uid=%s", username),
	)

	var bindErr error
	for _, bindDN := range bindDNs {
		bindErr = conn.Bind(bindDN, password)
		if bindErr == nil {
			fmt.Printf("✓ Authentication successful (bind DN: %s)\n", bindDN)

			// Perform a simple search to verify the connection is fully functional
			verifyRequest := ldap.NewSearchRequest(
				"",
				ldap.ScopeBaseObject,
				ldap.NeverDerefAliases,
				0,
				timeoutSec,
				false,
				"(objectClass=*)",
				[]string{"namingContexts"},
				nil,
			)

			verifyResult, err := conn.Search(verifyRequest)
			if err != nil {
				return fmt.Errorf("bind succeeded but search failed: %w", err)
			}

			if len(verifyResult.Entries) > 0 {
				fmt.Println("✓ Directory search successful")
			}
			return nil
		}
	}

	return fmt.Errorf("authentication failed: %w", bindErr)
}

// ldapDNToDomain converts an LDAP DN like "DC=office,DC=local" to "office.local"
func ldapDNToDomain(dn string) string {
	parts, err := ldap.ParseDN(dn)
	if err != nil || parts == nil {
		return ""
	}

	var domainParts []string
	for _, rdn := range parts.RDNs {
		for _, attr := range rdn.Attributes {
			if attr.Type == "DC" {
				domainParts = append(domainParts, attr.Value)
			}
		}
	}

	if len(domainParts) == 0 {
		return ""
	}

	domain := ""
	for _, part := range domainParts {
		if domain != "" {
			domain += "."
		}
		domain += part
	}
	return domain
}

// checkLDAPBatch checks multiple username/password pairs from a CSV file
func checkLDAPBatch(server, csvFile string, useTLS bool, timeoutSec int) (*BatchCheckResult, error) {
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

		err := checkLDAP(server, username, password, useTLS, timeoutSec)
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
