# Usage Guide

## Installation

```bash
go get github.com/yourorg/httpclient
```

## Import

```go
import "github.com/yourorg/httpclient"
```

## Basic Usage

### One-Line Requests

```go
// GET request
data, err := httpclient.GET("https://api.example.com/users")

// POST with JSON body
user := User{Name: "John", Email: "john@example.com"}
data, err := httpclient.POST("https://api.example.com/users", user)

// PUT request
data, err := httpclient.PUT("https://api.example.com/users/1", user)

// PATCH request
updates := map[string]string{"name": "Updated Name"}
data, err := httpclient.PATCH("https://api.example.com/users/1", updates)

// DELETE request
data, err := httpclient.DELETE("https://api.example.com/users/1")

// HEAD request (returns only error, no body)
err := httpclient.HEAD("https://api.example.com/users/1")

// OPTIONS request
data, err := httpclient.OPTIONS("https://api.example.com/users")
```

### JSON Requests (Recommended)

```go
// GET with automatic JSON parsing
var user User
err := httpclient.JSON("GET", "https://api.example.com/users/1", nil, &user)

// POST with JSON body and response parsing
newUser := User{Name: "John", Email: "john@example.com"}
var createdUser User
err := httpclient.JSON("POST", "https://api.example.com/users", newUser, &createdUser)

// PUT with JSON
var updatedUser User
err := httpclient.JSON("PUT", "https://api.example.com/users/1", user, &updatedUser)
```

## Custom Client Configuration

### Basic Configuration

```go
client := httpclient.New().
    WithTimeout(30 * time.Second).
    WithRetries(5).
    WithUserAgent("MyApp/1.0")

data, err := client.GET("https://api.example.com/users")
```

### Authentication

```go
// Bearer token authentication
client := httpclient.New().WithAuth("your-bearer-token")

// API key authentication
client := httpclient.New().WithAPIKey("X-API-Key", "your-api-key")

// Custom headers
client := httpclient.New().
    WithHeader("Authorization", "Custom auth-scheme token").
    WithHeader("X-Custom-Header", "custom-value")

// Multiple headers at once
headers := map[string]string{
    "X-API-Key":    "your-key",
    "X-Client-ID":  "your-client-id",
    "X-Version":    "v1",
}
client := httpclient.New().WithHeaders(headers)
```

### Base URL for API Clients

```go
client := httpclient.New().WithBaseURL("https://api.example.com")

// Now you can use relative URLs
var user User
err := client.JSON("GET", "/users/1", nil, &user)  // Requests https://api.example.com/users/1
```

## Advanced Features

### Rate Limiting

```go
// Limit to 100 requests per second
client := httpclient.New().WithRateLimiter(100)
```

### Circuit Breaker

```go
// Break circuit after 5 failures, recover after 60 seconds
client := httpclient.New().WithCircuitBreaker(5, 60*time.Second)
```

### Response Caching

```go
// Cache responses for 5 minutes
client := httpclient.New().WithCache(5 * time.Minute)
```

### Observability

```go
// Enable Prometheus metrics
client := httpclient.New().WithMetrics(true)

// Enable OpenTelemetry tracing
client := httpclient.New().WithTracing(true)

// Enable debug logging
client := httpclient.New().WithDebug(true)

// Enable all observability features
client := httpclient.New().
    WithMetrics(true).
    WithTracing(true).
    WithDebug(true)
```

### Combining Features

```go
// Production-ready client with advanced features
client := httpclient.New().
    WithBaseURL("https://api.example.com").
    WithAuth("your-token").
    WithTimeout(30 * time.Second).
    WithRetries(3).
    WithRateLimiter(50).
    WithCircuitBreaker(5, 60*time.Second).
    WithCache(5 * time.Minute).
    WithLoadBalancer(endpoints, "round-robin").
    WithHealthCheck(30*time.Second, "/health").
    WithCompression(true).
    WithBackupEndpoints(backups).
    WithMetrics(true).
    WithTracing(true).
    WithUserAgent("MyApp/1.0")
```

## Advanced Features

### Load Balancing

```go
// Round-robin load balancing
client := httpclient.New().
    WithLoadBalancer([]string{
        "https://api1.example.com",
        "https://api2.example.com",
        "https://api3.example.com",
    }, "round-robin")

// Random load balancing
client := httpclient.New().
    WithLoadBalancer(endpoints, "random")

// Least connections load balancing
client := httpclient.New().
    WithLoadBalancer(endpoints, "least-conn")
```

### Health Checks

```go
// Automatic health checking every 30 seconds
client := httpclient.New().
    WithHealthCheck(30*time.Second, "/health")

// Combined with load balancing
client := httpclient.New().
    WithLoadBalancer(endpoints, "round-robin").
    WithHealthCheck(60*time.Second, "/api/health")
```

### Request/Response Compression

```go
// Enable automatic compression
client := httpclient.New().
    WithCompression(true)

// Requests and responses are automatically compressed/decompressed
data, err := client.POST("https://api.example.com/data", largePayload)
```

### Security Features

```go
// Request signing with RSA
client := httpclient.New().
    WithRequestSigning("my-key-id", privateKeyPEM)

// IP whitelisting
client := httpclient.New().
    WithIPWhitelist([]string{
        "127.0.0.1",
        "192.168.1.0/24",
        "10.0.0.0/8",
    })

// Custom TLS configuration
client := httpclient.New().
    WithTLSConfig(&tls.Config{
        MinVersion: tls.VersionTLS13,
        CipherSuites: []uint16{
            tls.TLS_AES_256_GCM_SHA384,
        },
    })
```

### Request/Response Interceptors

```go
// Request interceptor
client := httpclient.New().
    WithRequestInterceptor(func(req *http.Request) error {
        req.Header.Set("X-Request-ID", generateRequestID())
        req.Header.Set("X-Timestamp", time.Now().Format(time.RFC3339))
        return nil
    })

// Response interceptor
client := httpclient.New().
    WithResponseInterceptor(func(resp *http.Response) error {
        if resp.Header.Get("X-Rate-Limit-Remaining") == "0" {
            log.Println("Rate limit exhausted")
        }
        return nil
    })

// Multiple interceptors
client := httpclient.New().
    WithRequestInterceptor(addAuthHeaders).
    WithRequestInterceptor(addTrackingHeaders).
    WithResponseInterceptor(logResponse).
    WithResponseInterceptor(checkRateLimit)
```

### Backup Endpoints and Failover

```go
// Automatic failover to backup endpoints
client := httpclient.New().
    WithBaseURL("https://primary-api.example.com").
    WithBackupEndpoints([]string{
        "https://backup1-api.example.com",
        "https://backup2-api.example.com",
    })

// If primary fails, automatically tries backups
data, err := client.GET("/users")
```

### Connection Management

```go
// Custom connection pool settings
client := httpclient.New().
    WithConnectionPool(100, 20).  // 100 max idle, 20 per host
    WithKeepAlive(60 * time.Second)

// Custom transport
transport := &http.Transport{
    MaxIdleConns: 200,
    // ... other settings
}
client := httpclient.New().
    WithCustomTransport(transport)
```

### Proxy Support

```go
// HTTP proxy
client := httpclient.New().
    WithProxy("http://proxy.example.com:8080")

// SOCKS proxy
client := httpclient.New().
    WithProxy("socks5://proxy.example.com:1080")
```

### Cookie Management

```go
// Automatic cookie handling
jar, _ := cookiejar.New(nil)
client := httpclient.New().
    WithCookieJar(jar)

// Cookies are automatically stored and sent
client.GET("https://example.com/login")  // Sets cookies
client.GET("https://example.com/profile") // Uses stored cookies
```

### Redirect Control

```go
// Custom redirect policy
client := httpclient.New().
    WithRedirectPolicy(func(req *http.Request, via []*http.Request) error {
        if len(via) >= 5 {
            return fmt.Errorf("too many redirects")
        }
        // Don't follow redirects to different domains
        if req.URL.Host != via[0].URL.Host {
            return fmt.Errorf("redirect to different domain not allowed")
        }
        return nil
    })
```

## Context-Aware Requests

### Using Context

```go
// With timeout context
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

data, err := httpclient.GetContext(ctx, "https://api.example.com/users")

// JSON with context
var user User
err := httpclient.JSONContext(ctx, "GET", "https://api.example.com/users/1", nil, &user)
```

### Custom Client with Context

```go
client := httpclient.New().WithTimeout(10 * time.Second)

ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

data, err := client.GetContext(ctx, "https://api.example.com/users")
```

## Real-World Examples

### API Client Service

```go
package api

import (
    "fmt"
    "github.com/yourorg/httpclient"
)

type Client struct {
    http httpclient.Client
}

func New(baseURL, token string) *Client {
    client := httpclient.New().
        WithBaseURL(baseURL).
        WithAuth(token).
        WithTimeout(30 * time.Second).
        WithRetries(3).
        WithLoadBalancer([]string{baseURL}, "round-robin").
        WithHealthCheck(60*time.Second, "/health").
        WithCompression(true).
        WithMetrics(true)

    return &Client{http: client}
}

func (c *Client) GetUser(id int) (*User, error) {
    var user User
    err := c.http.JSON("GET", fmt.Sprintf("/users/%d", id), nil, &user)
    if err != nil {
        return nil, fmt.Errorf("failed to get user %d: %w", id, err)
    }
    return &user, nil
}

func (c *Client) CreateUser(user *User) (*User, error) {
    var created User
    err := c.http.JSON("POST", "/users", user, &created)
    if err != nil {
        return nil, fmt.Errorf("failed to create user: %w", err)
    }
    return &created, nil
}

func (c *Client) UpdateUser(id int, user *User) (*User, error) {
    var updated User
    err := c.http.JSON("PUT", fmt.Sprintf("/users/%d", id), user, &updated)
    if err != nil {
        return nil, fmt.Errorf("failed to update user %d: %w", id, err)
    }
    return &updated, nil
}

func (c *Client) DeleteUser(id int) error {
    _, err := c.http.DELETE(fmt.Sprintf("/users/%d", id))
    if err != nil {
        return fmt.Errorf("failed to delete user %d: %w", id, err)
    }
    return nil
}
```

### Microservice Communication

```go
// High-availability microservice client
func NewUserService(endpoints []string, token string) *UserService {
    client := httpclient.New().
        WithLoadBalancer(endpoints, "least-conn").
        WithHealthCheck(30*time.Second, "/health").
        WithAuth(token).
        WithTimeout(10 * time.Second).
        WithRetries(5).
        WithCircuitBreaker(10, 120*time.Second).
        WithBackupEndpoints([]string{
            "https://backup-user-service.example.com",
        }).
        WithCompression(true).
        WithMetrics(true).
        WithTracing(true)

    return &UserService{client: client}
}

// Enterprise security configuration
func NewSecureClient(baseURL, keyID, privateKey string) httpclient.Client {
    return httpclient.New().
        WithBaseURL(baseURL).
        WithTLSConfig(&tls.Config{
            MinVersion: tls.VersionTLS13,
        }).
        WithRequestSigning(keyID, privateKey).
        WithIPWhitelist([]string{
            "10.0.0.0/8",
            "172.16.0.0/12",
            "192.168.0.0/16",
        }).
        WithRequestInterceptor(func(req *http.Request) error {
            req.Header.Set("X-Security-Level", "high")
            req.Header.Set("X-Timestamp", time.Now().Format(time.RFC3339))
            return nil
        }).
        WithTimeout(30 * time.Second).
        WithRetries(3)
}
```

### Performance-Optimized Client

```go
// High-performance client for heavy workloads
func NewHighPerformanceClient() httpclient.Client {
    client := httpclient.New().
        WithCompression(true).
        WithConnectionPool(200, 50).
        WithKeepAlive(120 * time.Second).
        WithCache(15 * time.Minute).
        WithRateLimiter(1000).
        WithTimeout(60 * time.Second).
        WithLoadBalancer([]string{
            "https://api1.example.com",
            "https://api2.example.com",
            "https://api3.example.com",
            "https://api4.example.com",
        }, "least-conn").
        WithHealthCheck(15*time.Second, "/health")

    return client
}
```

### Complete Enterprise Setup

```go
// Production-ready enterprise client
func NewEnterpriseClient(config EnterpriseConfig) httpclient.Client {
    return httpclient.New().
        WithBaseURL(config.BaseURL).
        WithAuth(config.Token).
        WithAPIKey("X-API-Key", config.APIKey).
        WithTimeout(config.Timeout).
        WithRetries(config.MaxRetries).
        WithLoadBalancer(config.Endpoints, "round-robin").
        WithHealthCheck(config.HealthCheckInterval, "/health").
        WithRateLimiter(config.RateLimit).
        WithCircuitBreaker(config.CircuitBreakerThreshold, config.CircuitBreakerTimeout).
        WithCache(config.CacheTTL).
        WithCompression(true).
        WithConnectionPool(config.MaxConnections, config.MaxConnectionsPerHost).
        WithKeepAlive(config.KeepAlive).
        WithTLSConfig(&tls.Config{
            MinVersion: tls.VersionTLS13,
        }).
        WithRequestSigning(config.SigningKeyID, config.PrivateKey).
        WithIPWhitelist(config.AllowedIPs).
        WithBackupEndpoints(config.BackupEndpoints).
        WithProxy(config.ProxyURL).
        WithMetrics(true).
        WithTracing(true).
        WithDebug(config.Debug).
        WithRequestInterceptor(func(req *http.Request) error {
            req.Header.Set("X-Client-Version", config.ClientVersion)
            req.Header.Set("X-Environment", config.Environment)
            return nil
        }).
        WithResponseInterceptor(func(resp *http.Response) error {
            // Log slow responses
            if resp.Header.Get("X-Response-Time") != "" {
                // Handle slow response logging
            }
            return nil
        })
}

type EnterpriseConfig struct {
    BaseURL                   string
    Token                     string
    APIKey                    string
    Timeout                   time.Duration
    MaxRetries                int
    Endpoints                 []string
    HealthCheckInterval       time.Duration
    RateLimit                 int
    CircuitBreakerThreshold   int
    CircuitBreakerTimeout     time.Duration
    CacheTTL                  time.Duration
    MaxConnections            int
    MaxConnectionsPerHost     int
    KeepAlive                 time.Duration
    SigningKeyID              string
    PrivateKey                string
    AllowedIPs                []string
    BackupEndpoints           []string
    ProxyURL                  string
    ClientVersion             string
    Environment               string
    Debug                     bool
}

// Order service client
type OrderService struct {
    client httpclient.Client
}

func NewOrderService(baseURL, token string) *OrderService {
    client := httpclient.New().
        WithBaseURL(baseURL).
        WithAuth(token).
        WithTimeout(10 * time.Second).
        WithRetries(5).
        WithCircuitBreaker(5, 60*time.Second).
        WithHeader("X-Service", "order-service").
        WithMetrics(true).
        WithTracing(true)

    return &OrderService{client: client}
}

func (s *OrderService) GetOrder(id int) (*Order, error) {
    var order Order
    err := s.client.JSON("GET", fmt.Sprintf("/orders/%d", id), nil, &order)
    return &order, err
}

func (s *OrderService) CreateOrder(order *Order) (*Order, error) {
    var created Order
    err := s.client.JSON("POST", "/orders", order, &created)
    return &created, err
}

// Business logic using both services
func ProcessOrder(userService *UserService, orderService *OrderService, orderID int) error {
    // Get order details
    order, err := orderService.GetOrder(orderID)
    if err != nil {
        return fmt.Errorf("failed to get order: %w", err)
    }

    // Get user details
    user, err := userService.GetUser(order.UserID)
    if err != nil {
        return fmt.Errorf("failed to get user: %w", err)
    }

    // Process order logic here...
    fmt.Printf("Processing order %d for user %s\n", order.ID, user.Name)
    
    return nil
}
```

### Configuration from Environment

```go
package config

import (
    "os"
    "strconv"
    "time"
    "github.com/yourorg/httpclient"
)

func NewHTTPClient() httpclient.Client {
    client := httpclient.New()

    // Configure from environment variables
    if baseURL := os.Getenv("API_BASE_URL"); baseURL != "" {
        client = client.WithBaseURL(baseURL)
    }

    if token := os.Getenv("API_TOKEN"); token != "" {
        client = client.WithAuth(token)
    }

    if timeoutStr := os.Getenv("API_TIMEOUT"); timeoutStr != "" {
        if timeout, err := time.ParseDuration(timeoutStr); err == nil {
            client = client.WithTimeout(timeout)
        }
    }

    if retriesStr := os.Getenv("API_RETRIES"); retriesStr != "" {
        if retries, err := strconv.Atoi(retriesStr); err == nil {
            client = client.WithRetries(retries)
        }
    }

    if os.Getenv("API_DEBUG") == "true" {
        client = client.WithDebug(true)
    }

    if os.Getenv("API_METRICS") == "true" {
        client = client.WithMetrics(true)
    }

    if os.Getenv("API_TRACING") == "true" {
        client = client.WithTracing(true)
    }

    return client
}
```

## Error Handling

### Understanding Errors

```go
data, err := httpclient.GET("https://api.example.com/users")
if err != nil {
    // Errors include HTTP status codes and response bodies
    fmt.Printf("Request failed: %v\n", err)
    // Output examples:
    // "HTTP 404: User not found"
    // "HTTP 500: Internal server error"
    // "request failed: context deadline exceeded"
}
```

### Handling Different Error Types

```go
data, err := httpclient.GET("https://api.example.com/users")
if err != nil {
    switch {
    case strings.Contains(err.Error(), "HTTP 404"):
        // Handle not found
        fmt.Println("Resource not found")
    case strings.Contains(err.Error(), "HTTP 5"):
        // Handle server errors
        fmt.Println("Server error, retrying...")
    case strings.Contains(err.Error(), "context deadline exceeded"):
        // Handle timeout
        fmt.Println("Request timed out")
    default:
        // Handle other errors
        fmt.Printf("Unexpected error: %v\n", err)
    }
}
```

## Testing

### Testing with httptest

```go
package api_test

import (
    "net/http"
    "net/http/httptest"
    "testing"
    "github.com/yourorg/httpclient"
)

func TestGetUser(t *testing.T) {
    // Create test server
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Method != "GET" {
            t.Errorf("Expected GET, got %s", r.Method)
        }
        if r.URL.Path != "/users/1" {
            t.Errorf("Expected /users/1, got %s", r.URL.Path)
        }
        
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"id": 1, "name": "John", "email": "john@example.com"}`))
    }))
    defer server.Close()

    // Test the client
    client := httpclient.New().WithBaseURL(server.URL)
    
    var user User
    err := client.JSON("GET", "/users/1", nil, &user)
    
    if err != nil {
        t.Fatalf("Request failed: %v", err)
    }
    
    if user.ID != 1 {
        t.Errorf("Expected ID 1, got %d", user.ID)
    }
    
    if user.Name != "John" {
        t.Errorf("Expected name John, got %s", user.Name)
    }
}
```

### Testing Error Conditions

```go
func TestErrorHandling(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusNotFound)
        w.Write([]byte("User not found"))
    }))
    defer server.Close()

    client := httpclient.New().WithBaseURL(server.URL)
    
    var user User
    err := client.JSON("GET", "/users/999", nil, &user)
    
    if err == nil {
        t.Fatal("Expected error for 404 response")
    }
    
    if !strings.Contains(err.Error(), "HTTP 404") {
        t.Errorf("Expected 404 error, got: %v", err)
    }
}
```

## Best Practices

### 1. Use JSON Methods for APIs

```go
// Preferred: Automatic JSON handling
var user User
err := httpclient.JSON("GET", "https://api.example.com/users/1", nil, &user)

// Avoid: Manual JSON handling
data, err := httpclient.GET("https://api.example.com/users/1")
if err != nil { return err }
err = json.Unmarshal(data, &user)
```

### 2. Configure Clients Once, Use Everywhere

```go
// Good: Create configured client once
var apiClient = httpclient.New().
    WithBaseURL("https://api.example.com").
    WithAuth("token").
    WithTimeout(30 * time.Second).
    WithRetries(3)

func GetUser(id int) (*User, error) {
    var user User
    err := apiClient.JSON("GET", fmt.Sprintf("/users/%d", id), nil, &user)
    return &user, err
}

// Avoid: Creating new clients repeatedly
func GetUser(id int) (*User, error) {
    client := httpclient.New().WithAuth("token")  // Don't do this
    // ...
}
```

### 3. Use Context for Timeouts

```go
// Good: Use context for request-specific timeouts
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

data, err := httpclient.GetContext(ctx, "https://api.example.com/users")

// Also good: Configure client timeout for all requests
client := httpclient.New().WithTimeout(30 * time.Second)
```

### 4. Handle Errors Appropriately

```go
// Good: Check and handle errors
data, err := httpclient.GET("https://api.example.com/users")
if err != nil {
    log.Printf("API request failed: %v", err)
    return fmt.Errorf("failed to fetch users: %w", err)
}

// Avoid: Ignoring errors
data, _ := httpclient.GET("https://api.example.com/users")  // Don't do this
```

### 5. Use Appropriate Features for Your Use Case

```go
// For high-traffic production services
client := httpclient.New().
    WithRateLimiter(100).
    WithCircuitBreaker(5, 60*time.Second).
    WithMetrics(true).
    WithTracing(true)

// For simple scripts or low-traffic applications
client := httpclient.New().
    WithTimeout(30 * time.Second).
    WithRetries(3)

// For development and debugging
client := httpclient.New().
    WithDebug(true)
```

This covers the most common usage patterns. The client is designed to be intuitive - if you need something that's not covered here, try the fluent interface approach and it will likely work as expected!