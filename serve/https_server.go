package serve

import (
	"crypto/subtle"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/UT-CTF/landschaft/util"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
)

// ServeDirectoryWithHTTPS serves the contents of a directory over HTTPS with basic auth
func ServeDirectoryWithHTTPS(dirPath string, port int) error {
	// Get certificate paths from existing function
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("error getting home directory: %v", err)
	}
	certDirectory = strings.Replace(certDirectory, "~", homeDir, 1)

	certFile, keyFile, err := getCertFiles(certDirectory)
	if err != nil {
		return fmt.Errorf("failed to get certificate files: %v", err)
	}

	// Check if directory exists
	dirInfo, err := os.Stat(dirPath)
	if err != nil {
		return fmt.Errorf("error accessing directory %s: %v", dirPath, err)
	}
	if !dirInfo.IsDir() {
		return fmt.Errorf("%s is not a directory", dirPath)
	}

	// Get absolute path for better logging
	absPath, err := filepath.Abs(dirPath)
	if err != nil {
		absPath = dirPath // Fallback to relative path
	}

	// Create a file server handler for the directory
	fileServer := http.FileServer(http.Dir(dirPath))

	// generate random password
	defaultPassword := util.GenerateRandomDefaultPassword(12)
	var password string

	// prompt for password
	huh.NewInput().
		Title("Enter password for default auth (enter for random default):").
		Prompt("> ").
		Placeholder(defaultPassword).
		Value(&password).Run()

	if password == "" {
		password = defaultPassword
	}

	// Wrap with basic auth middleware
	authHandler := basicAuthMiddleware(fileServer, "admin", password)

	// Register the handler
	http.Handle("/", authHandler)

	// Configure TLS with modern cipher suites and TLS versions
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}

	// Create server with TLS config
	server := &http.Server{
		Addr:      fmt.Sprintf(":%d", port),
		TLSConfig: tlsConfig,
	}

	// Log startup information
	log.Info("Starting HTTPS server for", "directory", absPath)
	log.Infof("Server URL: %s", util.PurpleStyle.Render(fmt.Sprintf("https://localhost:%d/", port)))
	log.Info("Basic auth credentials:", "username", "admin", "password", password)

	// Start the HTTPS server
	return server.ListenAndServeTLS(certFile, keyFile)
}

// basicAuthMiddleware wraps a handler with HTTP Basic Authentication
func basicAuthMiddleware(next http.Handler, username, password string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get credentials from request
		user, pass, ok := r.BasicAuth()

		// Check if credentials match using constant time comparison to prevent timing attacks
		userMatch := subtle.ConstantTimeCompare([]byte(user), []byte(username)) == 1
		passMatch := subtle.ConstantTimeCompare([]byte(pass), []byte(password)) == 1

		// If authentication fails, request authentication
		if !ok || !userMatch || !passMatch {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// If authentication succeeds, call the next handler
		next.ServeHTTP(w, r)
	})
}
