package serve

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"path"
	"time"

	"github.com/charmbracelet/log"
)

// raaaaa we need just one cert generation function, not different files
// todo: merge cert generation

// todo: merge this out into general utils
func getBindingAddresses() ([]net.IP, error) {
	// get all network interfaces
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to get network interfaces: %v", err)
	}
	var addresses []net.IP
	for _, iface := range interfaces {
		// get all addresses for the interface
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, fmt.Errorf("failed to get addresses for interface %s: %v", iface.Name, err)
		}
		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok && ipNet.IP.To4() != nil {
				addresses = append(addresses, ipNet.IP)
			}
		}
	}
	if len(addresses) == 0 {
		return nil, fmt.Errorf("no valid IP addresses found")
	}
	return addresses, nil
}

func getHostname() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", fmt.Errorf("failed to get hostname: %v", err)
	}
	return hostname, nil
}

func generateCert() ([]byte, []byte, error) {
	// Generate a private key
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate private key: %v", err)
	}

	// Create certificate template
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate serial number: %v", err)
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(365 * 24 * time.Hour) // Valid for 1 year

	hostname, err := getHostname()
	if err != nil {
		log.Error("Error getting hostname, using localhost")
		hostname = "localhost"
	}

	dnsNames := []string{hostname}
	addresses, _ := getBindingAddresses()

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Landschaft"},
			CommonName:   hostname,
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IPAddresses:           addresses,
		DNSNames:              dnsNames,
	}

	// Create the certificate
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create certificate: %v", err)
	}

	// Encode certificate to PEM
	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: derBytes,
	})

	// Encode private key to PEM
	privBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal private key: %v", err)
	}

	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privBytes,
	})

	return certPEM, keyPEM, nil
}

func getCertFiles(directory string) (string, string, error) {
	certPath := path.Join(directory, "cert.pem")
	keyPath := path.Join(directory, "key.pem")

	// Check if both the certificate and key files already exist, if so, return their paths
	if _, err := os.Stat(certPath); err == nil {
		return certPath, keyPath, nil
	}
	if _, err := os.Stat(keyPath); err == nil {
		return certPath, keyPath, nil
	}
	// If not, generate new certificate and key
	certPEM, keyPEM, err := generateCert()

	// make sure directory exists
	if err := os.MkdirAll(directory, 0755); err != nil {
		return "", "", fmt.Errorf("failed to create directory: %v", err)
	}

	if err != nil {
		return "", "", fmt.Errorf("failed to generate certificate: %v", err)
	}
	// Write the certificate and key to files
	if err := os.WriteFile(certPath, certPEM, 0644); err != nil {
		return "", "", fmt.Errorf("failed to write certificate file: %v", err)
	}
	if err := os.WriteFile(keyPath, keyPEM, 0600); err != nil {
		return "", "", fmt.Errorf("failed to write key file: %v", err)
	}

	return certPath, keyPath, nil

}
