package tcvectordb

import (
	"crypto/tls"
	"crypto/x509"
	"net"
	"os"
	"strings"

	"github.com/pkg/errors"
)

// configureTLSCertificates configures TLS certificates based on the provided CA certificate option
func configureTLSCertificates(tlsConfig *tls.Config, caCertOption string) error {
	if caCertOption == "" {
		return nil
	}

	caCertPool := x509.NewCertPool()
	certContent := strings.TrimSpace(caCertOption)

	// Check if CACert is a file path
	if _, err := os.Stat(certContent); err == nil {
		// It's a file path, read the certificate from file
		certBytes, err := os.ReadFile(certContent)
		if err != nil {
			return errors.Wrapf(err, "failed to read CA certificate file: %s", certContent)
		}
		certContent = string(certBytes)
	}

	if ok := caCertPool.AppendCertsFromPEM([]byte(certContent)); !ok {
		return errors.Errorf("failed to parse CA certificate: invalid PEM format or certificate content")
	}
	tlsConfig.RootCAs = caCertPool
	return nil
}

// CreateTLSConfig creates a complete TLS configuration based on client options and URL
func CreateTLSConfig(option *ClientOption, url string) (*tls.Config, error) {
	if option.CACert == "" && !option.InsecureSkipVerify {
		return nil, nil
	}
	tlsConfig := &tls.Config{
		InsecureSkipVerify: option.InsecureSkipVerify,
	}
	if option.InsecureSkipVerify {
		return tlsConfig, nil
	}

	// Configure CA certificate if provided
	if option.CACert != "" {
		if err := configureTLSCertificates(tlsConfig, option.CACert); err != nil {
			return nil, err
		}
	}

	// If URL is an IP address, set ServerName to fixed value
	if isIPAddress(url) {
		tlsConfig.ServerName = "vdb.tencentcloud.com"
	}
	// If URL is not an IP and ServerName is not specified, use default behavior

	return tlsConfig, nil
}

// isIPAddress checks if the given string is an IP address
func isIPAddress(s string) bool {
	// Remove protocol prefix if present
	host := strings.TrimPrefix(strings.TrimPrefix(s, "http://"), "https://")

	// Remove port if present
	if colonIndex := strings.LastIndex(host, ":"); colonIndex != -1 {
		host = host[:colonIndex]
	}

	// Use net.ParseIP for accurate IP address detection
	return net.ParseIP(host) != nil
}
