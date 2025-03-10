package graylog

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"time"
)

func genCa() error {
	// Define file paths for CA certificate and private key
	certPath := "ca.crt"
	keyPath := "ca.key"
	bundlePath := "ca-bundle.crt"

	// Check if files already exist to prevent overwriting
	if _, err := os.Stat(certPath); err == nil {
		return fmt.Errorf("CA certificate file already exists: %s", certPath)
	}
	if _, err := os.Stat(keyPath); err == nil {
		return fmt.Errorf("CA private key file already exists: %s", keyPath)
	}
	if _, err := os.Stat(bundlePath); err == nil {
		return fmt.Errorf("CA bundle file already exists: %s", bundlePath)
	}

	// Generate RSA private key for CA
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return fmt.Errorf("failed to generate private key: %w", err)
	}

	// Create certificate template for CA
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Landschaft CA"},
			CommonName:   "Landschaft Root CA",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0), // Valid for 10 years
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLen:            1,
	}

	// Create self-signed certificate
	certBytes, err := x509.CreateCertificate(
		rand.Reader,
		&template,
		&template,
		&privateKey.PublicKey,
		privateKey,
	)
	if err != nil {
		return fmt.Errorf("failed to create certificate: %w", err)
	}

	// Write private key to file (O_EXCL ensures we don't overwrite)
	keyFile, err := os.OpenFile(keyPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return fmt.Errorf("failed to create key file: %w", err)
	}
	defer keyFile.Close()

	pkcs8Bytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return fmt.Errorf("failed to marshal private key: %w", err)
	}

	err = pem.Encode(keyFile, &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: pkcs8Bytes,
	})
	if err != nil {
		return fmt.Errorf("failed to write private key: %w", err)
	}

	// Write certificate to file (O_EXCL ensures we don't overwrite)
	certFile, err := os.OpenFile(certPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return fmt.Errorf("failed to create certificate file: %w", err)
	}
	defer certFile.Close()

	certBlock := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	}

	err = pem.Encode(certFile, certBlock)
	if err != nil {
		return fmt.Errorf("failed to write certificate: %w", err)
	}

	// Write bundle file with both certificate and private key
	bundleFile, err := os.OpenFile(bundlePath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return fmt.Errorf("failed to create bundle file: %w", err)
	}
	defer bundleFile.Close()

	// Write certificate to bundle first
	err = pem.Encode(bundleFile, certBlock)
	if err != nil {
		return fmt.Errorf("failed to write certificate to bundle file: %w", err)
	}

	// Then write private key to bundle
	privateKeyBlock := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: pkcs8Bytes,
	}

	err = pem.Encode(bundleFile, privateKeyBlock)
	if err != nil {
		return fmt.Errorf("failed to write private key to bundle file: %w", err)
	}

	return nil
}
