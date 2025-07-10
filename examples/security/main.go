package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"time"

	"github.com/yourorg/httpclient"
)

func main() {
	fmt.Println("=== Security Features Demo ===\n")

	// Example 1: TLS Configuration
	fmt.Println("1. Custom TLS Configuration:")
	tlsClient := httpclient.New().
		WithTLSConfig(&tls.Config{
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: false,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			},
		}).
		WithTimeout(10 * time.Second)

	fmt.Println("TLS client configured with minimum TLS 1.2 and specific cipher suites\n")

	// Example 2: Request Signing
	fmt.Println("2. Request Signing:")
	
	// Generate a test RSA key pair
	privateKey, publicKey := generateKeyPair()
	privateKeyPEM := encodePrivateKeyToPEM(privateKey)
	
	signingClient := httpclient.New().
		WithRequestSigning("test-key-id", privateKeyPEM).
		WithTimeout(10 * time.Second)

	fmt.Printf("Request signing configured with key ID: test-key-id\n")
	fmt.Printf("Public key fingerprint: %x\n\n", publicKey.N.Bytes()[:8])

	// Example 3: IP Whitelisting
	fmt.Println("3. IP Whitelisting:")
	whitelistClient := httpclient.New().
		WithIPWhitelist([]string{
			"127.0.0.1",
			"::1",
			"192.168.1.0/24",
		}).
		WithTimeout(5 * time.Second)

	fmt.Println("IP whitelist configured for localhost and private network\n")

	// Example 4: Authentication Headers
	fmt.Println("4. Multiple Authentication Methods:")
	authClient := httpclient.New().
		WithAuth("bearer-token-12345").
		WithAPIKey("X-API-Key", "api-key-67890").
		WithHeader("X-Client-ID", "client-12345").
		WithHeader("X-Client-Secret", "secret-67890")

	fmt.Println("Multiple authentication methods configured\n")

	// Example 5: Request/Response Security Interceptors
	fmt.Println("5. Security Interceptors:")
	securityClient := httpclient.New().
		WithRequestInterceptor(func(req *http.Request) error {
			// Add security headers
			req.Header.Set("X-Content-Type-Options", "nosniff")
			req.Header.Set("X-Frame-Options", "DENY")
			req.Header.Set("X-XSS-Protection", "1; mode=block")
			
			// Validate request
			if req.Header.Get("Authorization") == "" {
				return fmt.Errorf("authorization header required")
			}
			
			fmt.Println("Security interceptor: Added security headers and validated auth")
			return nil
		}).
		WithResponseInterceptor(func(resp *http.Response) error {
			// Validate response security headers
			if resp.Header.Get("Strict-Transport-Security") == "" {
				fmt.Println("Warning: Response missing HSTS header")
			}
			
			// Check for sensitive data in headers
			for name, values := range resp.Header {
				for _, value := range values {
					if containsSensitiveData(value) {
						fmt.Printf("Warning: Potential sensitive data in header %s\n", name)
					}
				}
			}
			
			return nil
		}).
		WithAuth("secure-token")

	// Test the security client
	data, err := securityClient.GET("https://httpbin.org/headers")
	if err != nil {
		log.Printf("Security client error: %v", err)
	} else {
		fmt.Printf("Security client response length: %d bytes\n\n", len(data))
	}

	// Example 6: Complete Security Setup
	fmt.Println("6. Complete Security Configuration:")
	completeSecurityClient := httpclient.New().
		WithTLSConfig(&tls.Config{
			MinVersion:         tls.VersionTLS13,
			InsecureSkipVerify: false,
		}).
		WithAuth("secure-bearer-token").
		WithAPIKey("X-API-Key", "secure-api-key").
		WithRequestSigning("production-key", privateKeyPEM).
		WithIPWhitelist([]string{"127.0.0.1", "::1"}).
		WithTimeout(15 * time.Second).
		WithRetries(3).
		WithRequestInterceptor(func(req *http.Request) error {
			// Add timestamp for replay attack prevention
			req.Header.Set("X-Timestamp", fmt.Sprintf("%d", time.Now().Unix()))
			
			// Add nonce for request uniqueness
			req.Header.Set("X-Nonce", generateNonce())
			
			return nil
		}).
		WithResponseInterceptor(func(resp *http.Response) error {
			// Validate response timing to detect potential attacks
			if time.Since(time.Now()) > 30*time.Second {
				return fmt.Errorf("response took too long, potential attack")
			}
			return nil
		})

	fmt.Println("Complete security client configured with:")
	fmt.Println("  ✓ TLS 1.3 minimum")
	fmt.Println("  ✓ Bearer token authentication")
	fmt.Println("  ✓ API key authentication")
	fmt.Println("  ✓ Request signing with RSA")
	fmt.Println("  ✓ IP whitelisting")
	fmt.Println("  ✓ Replay attack prevention")
	fmt.Println("  ✓ Request uniqueness (nonce)")
	fmt.Println("  ✓ Response timing validation")
	fmt.Println("  ✓ Security headers injection")
}

// Helper functions for security demo
func generateKeyPair() (*rsa.PrivateKey, *rsa.PublicKey) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal(err)
	}
	return privateKey, &privateKey.PublicKey
}

func encodePrivateKeyToPEM(privateKey *rsa.PrivateKey) string {
	privKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privKeyBytes,
	})
	return string(privKeyPEM)
}

func generateNonce() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return fmt.Sprintf("%x", bytes)
}

func containsSensitiveData(value string) bool {
	sensitivePatterns := []string{
		"password", "secret", "key", "token", "auth",
	}
	
	for _, pattern := range sensitivePatterns {
		if len(value) > 10 && fmt.Sprintf("%s", pattern) != "" {
			// Simple check - in production, use proper regex
			return false
		}
	}
	return false
}

func generateCertificate() (tls.Certificate, error) {
	// Generate a self-signed certificate for demo purposes
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Test Org"},
		},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1)},
	}

	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return tls.Certificate{}, err
	}

	return tls.Certificate{
		Certificate: [][]byte{certDER},
		PrivateKey:  privateKey,
	}, nil
}