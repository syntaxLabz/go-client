package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/yourorg/httpclient"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	fmt.Println("=== Enterprise HTTP Client Features ===\n")

	// Example 1: Load Balancing with Health Checks
	fmt.Println("1. Load Balancing with Health Checks:")
	lbClient := httpclient.New().
		WithLoadBalancer([]string{
			"https://jsonplaceholder.typicode.com",
			"https://httpbin.org",
		}, "round-robin").
		WithHealthCheck(30*time.Second, "/health").
		WithTimeout(10 * time.Second)

	fmt.Println("Load balancer configured with health checks\n")

	// Example 2: Request/Response Compression
	fmt.Println("2. Compression and Custom Transport:")
	compressedClient := httpclient.New().
		WithCompression(true).
		WithConnectionPool(50, 10).
		WithKeepAlive(30 * time.Second)

	var user User
	err := compressedClient.JSON("GET", "https://jsonplaceholder.typicode.com/users/1", nil, &user)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("User with compression: %+v\n\n", user)
	}

	// Example 3: Request Signing and Security
	fmt.Println("3. Request Signing and IP Whitelisting:")
	secureClient := httpclient.New().
		WithIPWhitelist([]string{"127.0.0.1", "::1"}).
		WithTLSConfig(&tls.Config{
			MinVersion: tls.VersionTLS12,
		})

	fmt.Println("Secure client configured with IP whitelist and TLS settings\n")

	// Example 4: Cookie Jar and Redirect Policy
	fmt.Println("4. Cookie Management and Redirect Policy:")
	jar, _ := cookiejar.New(nil)
	cookieClient := httpclient.New().
		WithCookieJar(jar).
		WithRedirectPolicy(func(req *http.Request, via []*http.Request) error {
			if len(via) >= 3 {
				return fmt.Errorf("too many redirects")
			}
			return nil
		})

	fmt.Println("Cookie client configured with custom redirect policy\n")

	// Example 5: Request/Response Interceptors
	fmt.Println("5. Request/Response Interceptors:")
	interceptorClient := httpclient.New().
		WithRequestInterceptor(func(req *http.Request) error {
			req.Header.Set("X-Request-ID", "12345")
			fmt.Printf("Request interceptor: Added request ID\n")
			return nil
		}).
		WithResponseInterceptor(func(resp *http.Response) error {
			fmt.Printf("Response interceptor: Status %d\n", resp.StatusCode)
			return nil
		})

	data, err := interceptorClient.GET("https://jsonplaceholder.typicode.com/users/1")
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("Response length: %d bytes\n\n", len(data))
	}

	// Example 6: Backup Endpoints and Failover
	fmt.Println("6. Backup Endpoints and Failover:")
	failoverClient := httpclient.New().
		WithBaseURL("https://primary-api.example.com").
		WithBackupEndpoints([]string{
			"https://backup1-api.example.com",
			"https://backup2-api.example.com",
		}).
		WithTimeout(5 * time.Second).
		WithRetries(2)

	fmt.Println("Failover client configured with backup endpoints\n")

	// Example 7: Proxy Support
	fmt.Println("7. Proxy Configuration:")
	proxyClient := httpclient.New().
		WithProxy("http://proxy.example.com:8080").
		WithTimeout(15 * time.Second)

	fmt.Println("Proxy client configured\n")

	// Example 8: Complete Enterprise Setup
	fmt.Println("8. Complete Enterprise Configuration:")
	enterpriseClient := httpclient.New().
		WithBaseURL("https://api.enterprise.com").
		WithAuth("enterprise-token").
		WithTimeout(30 * time.Second).
		WithRetries(5).
		WithLoadBalancer([]string{
			"https://api1.enterprise.com",
			"https://api2.enterprise.com",
			"https://api3.enterprise.com",
		}, "least-conn").
		WithHealthCheck(60*time.Second, "/health").
		WithRateLimiter(200).
		WithCircuitBreaker(10, 120*time.Second).
		WithCache(10 * time.Minute).
		WithCompression(true).
		WithMetrics(true).
		WithTracing(true).
		WithConnectionPool(100, 20).
		WithKeepAlive(60 * time.Second).
		WithTLSConfig(&tls.Config{
			MinVersion: tls.VersionTLS13,
		}).
		WithRequestInterceptor(func(req *http.Request) error {
			req.Header.Set("X-Enterprise-Client", "v1.0")
			return nil
		}).
		WithBackupEndpoints([]string{
			"https://backup-api.enterprise.com",
		})

	fmt.Println("Enterprise client configured with all advanced features:")
	fmt.Println("  ✓ Load balancing (least-connection)")
	fmt.Println("  ✓ Health checks every 60s")
	fmt.Println("  ✓ Rate limiting (200 RPS)")
	fmt.Println("  ✓ Circuit breaker (10 failures, 2min timeout)")
	fmt.Println("  ✓ Response caching (10min TTL)")
	fmt.Println("  ✓ Request/response compression")
	fmt.Println("  ✓ Prometheus metrics")
	fmt.Println("  ✓ OpenTelemetry tracing")
	fmt.Println("  ✓ Connection pooling (100 max, 20 per host)")
	fmt.Println("  ✓ Keep-alive connections (60s)")
	fmt.Println("  ✓ TLS 1.3 minimum")
	fmt.Println("  ✓ Request interceptors")
	fmt.Println("  ✓ Backup endpoints for failover")
	fmt.Println("  ✓ Exponential backoff retries (5 attempts)")
	fmt.Println("  ✓ 30s request timeout")
}