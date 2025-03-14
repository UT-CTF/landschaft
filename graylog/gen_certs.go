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

func genCerts() error {
	if err := genCaCerts(); err != nil {
		return err
	}
	return genServerCerts("graylog.internal")
}

func genCaCerts() error {
	// Define file paths for CA certificate and private key
	certPath := "ca.crt"
	keyPath := "ca.key"
	bundlePath := "ca-bundle.key"

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

func genServerCerts(hostname string) error {
	// Define file paths
	caKeyPath := "ca.key"
	caCertPath := "ca.crt"
	serverKeyPath := hostname + ".key"
	serverCertPath := hostname + ".crt"
	serverBundlePath := hostname + ".bundle.crt"

	// Check if files already exist to prevent overwriting
	if _, err := os.Stat(serverCertPath); err == nil {
		return fmt.Errorf("server certificate file already exists: %s", serverCertPath)
	}
	if _, err := os.Stat(serverKeyPath); err == nil {
		return fmt.Errorf("server private key file already exists: %s", serverKeyPath)
	}
	if _, err := os.Stat(serverBundlePath); err == nil {
		return fmt.Errorf("server bundle file already exists: %s", serverBundlePath)
	}

	// Read the CA certificate
	caCertPEM, err := os.ReadFile(caCertPath)
	if err != nil {
		return fmt.Errorf("failed to read CA certificate: %w", err)
	}
	caCertBlock, _ := pem.Decode(caCertPEM)
	if caCertBlock == nil {
		return fmt.Errorf("failed to parse CA certificate PEM")
	}
	caCert, err := x509.ParseCertificate(caCertBlock.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse CA certificate: %w", err)
	}

	// Read the CA private key
	caKeyPEM, err := os.ReadFile(caKeyPath)
	if err != nil {
		return fmt.Errorf("failed to read CA private key: %w", err)
	}
	caKeyBlock, _ := pem.Decode(caKeyPEM)
	if caKeyBlock == nil {
		return fmt.Errorf("failed to parse CA private key PEM")
	}
	caKey, err := x509.ParsePKCS8PrivateKey(caKeyBlock.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse CA private key: %w", err)
	}
	caPrivKey, ok := caKey.(*rsa.PrivateKey)
	if !ok {
		return fmt.Errorf("CA private key is not RSA")
	}

	// Generate server private key
	serverPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return fmt.Errorf("failed to generate server private key: %w", err)
	}

	// Create certificate template for server
	serialNumber, err := rand.Int(rand.Reader, big.NewInt(1<<62))
	if err != nil {
		return fmt.Errorf("failed to generate serial number: %w", err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Landschaft"},
			CommonName:   hostname,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0), // Valid for 1 year
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  false,
		DNSNames:              []string{hostname},
	}

	// Create server certificate signed by CA
	certBytes, err := x509.CreateCertificate(
		rand.Reader,
		&template,
		caCert,
		&serverPrivKey.PublicKey,
		caPrivKey,
	)
	if err != nil {
		return fmt.Errorf("failed to create server certificate: %w", err)
	}

	// Write server private key to file
	keyFile, err := os.OpenFile(serverKeyPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return fmt.Errorf("failed to create server key file: %w", err)
	}
	defer keyFile.Close()

	pkcs8Bytes, err := x509.MarshalPKCS8PrivateKey(serverPrivKey)
	if err != nil {
		return fmt.Errorf("failed to marshal server private key: %w", err)
	}

	err = pem.Encode(keyFile, &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: pkcs8Bytes,
	})
	if err != nil {
		return fmt.Errorf("failed to write server private key: %w", err)
	}

	// Write server certificate to file
	certFile, err := os.OpenFile(serverCertPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return fmt.Errorf("failed to create server certificate file: %w", err)
	}
	defer certFile.Close()

	certBlock := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	}

	err = pem.Encode(certFile, certBlock)
	if err != nil {
		return fmt.Errorf("failed to write server certificate: %w", err)
	}

	// Write bundle file with both server certificate and CA certificate
	bundleFile, err := os.OpenFile(serverBundlePath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return fmt.Errorf("failed to create server bundle file: %w", err)
	}
	defer bundleFile.Close()

	// Write server certificate to bundle first
	err = pem.Encode(bundleFile, certBlock)
	if err != nil {
		return fmt.Errorf("failed to write server certificate to bundle file: %w", err)
	}

	// Then write CA certificate to bundle
	err = pem.Encode(bundleFile, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caCertBlock.Bytes,
	})
	if err != nil {
		return fmt.Errorf("failed to write CA certificate to bundle file: %w", err)
	}

	return nil
}
