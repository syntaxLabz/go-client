# Ultra-Simple HTTP Client for Go

The most developer-friendly HTTP client for Go. Production-ready with enterprise features, but feels like magic to use.

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/yourorg/httpclient)](https://goreportcard.com/report/github.com/yourorg/httpclient)

## Why This Client?

```go
// Before: Standard library
req, err := http.NewRequest("GET", "https://api.example.com/users", nil)
if err != nil { /* handle error */ }
req.Header.Set("Authorization", "Bearer token")
client := &http.Client{Timeout: 30 * time.Second}
resp, err := client.Do(req)
if err != nil { /* handle error */ }
defer resp.Body.Close()
body, err := io.ReadAll(resp.Body)
if err != nil { /* handle error */ }
var user User
json.Unmarshal(body, &user)

// After: This client
var user User
err := httpclient.New().WithAuth("token").JSON("GET", "https://api.example.com/users", nil, &user)
```

## Installation

```bash
go get github.com/syntaxLabz/go-clien
```

## Quick Start

### One-Line Requests

```go
import "github.com/yourorg/httpclient"

// GET request
data, err := httpclient.GET("https://api.example.com/users")

// POST with JSON
user := User{Name: "John"}
data, err := httpclient.POST("https://api.example.com/users", user)

// All HTTP methods supported
data, err := httpclient.PUT("https://api.example.com/users/1", user)
data, err := httpclient.PATCH("https://api.example.com/users/1", updates)
data, err := httpclient.DELETE("https://api.example.com/users/1")
err := httpclient.HEAD("https://api.example.com/users/1")
```

### JSON Made Effortless

```go
// GET and parse JSON in one line
var user User
err := httpclient.JSON("GET", "https://api.example.com/users/1", nil, &user)

// POST and parse response
var result Response
err := httpclient.JSON("POST", "https://api.example.com/users", newUser, &result)
```

### Smart Constructors for Different Use Cases

```go
// Microservices client with AI features
client := httpclient.NewForMicroservices()

// API client optimized for REST APIs
client := httpclient.NewForAPI()

// Enterprise client with full security
client := httpclient.NewForEnterprise()

// Development client with debugging tools
client := httpclient.NewForDevelopment()
```

### Advanced Features in One Line

```go
// Batch requests
responses, err := httpclient.Batch().
    Add("GET", "https://api.example.com/users/1", nil).
    Add("GET", "https://api.example.com/users/2", nil).
    Execute()

// WebSocket connection
ws, err := httpclient.WebSocket("wss://api.example.com/ws")

// GraphQL query
err := httpclient.GraphQL(query, variables, &result)

// HTTP streaming
stream, err := httpclient.Stream("GET", "https://api.example.com/stream", nil)
```
### Custom Configuration (Fluent Interface)

```go
client := httpclient.New().
    WithTimeout(30 * time.Second).
    WithRetries(5).
    WithAuth("your-bearer-token").
    WithBaseURL("https://api.example.com").
    WithHeader("X-API-Key", "your-key").
    WithUserAgent("MyApp/1.0").
    WithDebug(true)

var user User
err := client.JSON("GET", "/users/1", nil, &user)
```

## Advanced Features

### 🧠 AI-Powered Features

```go
client := httpclient.New().
    WithAIRetry(true).                    // Machine learning retry strategy
    WithSmartCaching(true).               // Intelligent caching decisions
    WithPredictivePreloading(true).       // Anticipate future requests
    WithAdaptiveTimeout(true).            // Dynamic timeout optimization
    WithPerformanceOptimization(true)     // Real-time performance tuning
```

### 🚀 Real-time & Streaming

```go
// HTTP Streaming
stream, err := client.Stream("GET", "https://api.example.com/events", nil)
for data := range stream {
    fmt.Printf("Received: %s\n", data)
}

// WebSocket
ws, err := client.WebSocket("wss://api.example.com/ws")
ws.Send("Hello!")
data, err := ws.Receive()

// Server-Sent Events
events, err := client.SSE("https://api.example.com/events")
for event := range events {
    fmt.Printf("Event: %s\n", event.Data)
}
```

### ⚡ Batch & Pipeline Operations

```go
// Concurrent batch execution
responses, err := client.Batch().
    Add("GET", "/users/1", nil).
    Add("GET", "/users/2", nil).
    Add("POST", "/users", newUser).
    Execute()

// Sequential pipeline with streaming results
pipeline, err := client.Pipeline().
    Add("GET", "/users/1", nil).
    Add("GET", "/posts/1", nil).
    Execute()

for response := range pipeline {
    fmt.Printf("Response %d: %v\n", response.Index, response.Data)
}
```

### 🎯 GraphQL Support

```go
query := `
    query GetUser($id: ID!) {
        user(id: $id) {
            name
            email
            posts {
                title
            }
        }
    }
`

variables := map[string]interface{}{"id": "123"}
var result UserResponse
err := client.GraphQL(query, variables, &result)
```
### Enterprise-Grade Features

```go
client := httpclient.New().
    WithRateLimiter(100).                          // 100 requests per second
    WithCircuitBreaker(5, 60*time.Second).        // Break after 5 failures for 60s
    WithCache(5 * time.Minute).                   // Cache responses for 5 minutes
    WithMetrics(true).                            // Prometheus metrics
    WithTracing(true).                            // OpenTelemetry tracing
    WithRetries(3).                               // Exponential backoff retries
    WithLoadBalancer(endpoints, "round-robin").   // Load balancing
    WithHealthCheck(30*time.Second, "/health").   // Health checks
    WithCompression(true).                        // Request/response compression
    WithRequestSigning("key-id", privateKey).     // Request signing
    WithIPWhitelist([]string{"127.0.0.1"}).      // IP whitelisting
    WithBackupEndpoints(backups)                  // Automatic failover
```

### Load Balancing & High Availability

```go
// Multiple load balancing strategies
client := httpclient.New().
    WithLoadBalancer([]string{
        "https://api1.example.com",
        "https://api2.example.com",
        "https://api3.example.com",
    }, "round-robin"). // or "random", "least-conn", "weighted-random"
    WithHealthCheck(30*time.Second, "/health").
    WithBackupEndpoints([]string{
        "https://backup-api.example.com",
    })
```

### Security Features

```go
// Comprehensive security configuration
client := httpclient.New().
    WithTLSConfig(&tls.Config{MinVersion: tls.VersionTLS13}).
    WithRequestSigning("key-id", privateKey).
    WithIPWhitelist([]string{"192.168.1.0/24"}).
    WithRequestInterceptor(func(req *http.Request) error {
        req.Header.Set("X-Security-Token", "secure-token")
        return nil
    })
```

### 🔧 Developer Experience Features

```go
// Development client with debugging tools
devClient := httpclient.NewForDevelopment().
    WithMocking(true).                    // Built-in mocking
    WithRecording(true).                  // Record/replay requests
    WithValidation(schema).               // Request/response validation
    WithChaosEngineering(chaosConfig).    // Fault injection testing
    WithDebug(true)                       // Comprehensive debugging

// Auto-retry with intelligent conditions
client := httpclient.New().
    WithAutoRetry(httpclient.AutoRetryConfig{
        MaxAttempts:     5,
        BackoffStrategy: "exponential",
        RetryConditions: []string{"timeout", "5xx", "connection_error"},
        JitterEnabled:   true,
    })
```

### 🌐 Advanced Networking

```go
client := httpclient.New().
    WithHTTP3(true).                      // HTTP/3 support
    WithMultipath(true).                  // Multipath TCP
    WithDNSOverHTTPS(true).              // DNS over HTTPS
    WithEdgeOptimization(true).           // CDN edge optimization
    WithProxy("socks5://proxy:1080")      // Advanced proxy support
```
### Performance Optimization

```go
// High-performance configuration
client := httpclient.New().
    WithCompression(true).
    WithConnectionPool(100, 20).
    WithKeepAlive(60 * time.Second).
    WithCache(10 * time.Minute).
    WithProxy("http://proxy.example.com:8080")
```

### Context Support

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

data, err := httpclient.GetContext(ctx, "https://api.example.com/users")
err := httpclient.JSONContext(ctx, "POST", "https://api.example.com/users", user, &result)
```

### Microservice Communication

```go
// User service client
userAPI := httpclient.New().
    WithBaseURL("http://user-service:8080").
    WithAuth("service-token").
    WithTimeout(5 * time.Second).
    WithRetries(3).
    WithMetrics(true)

// Order service client
orderAPI := httpclient.New().
    WithBaseURL("http://order-service:8080").
    WithAuth("service-token").
    WithTimeout(10 * time.Second).
    WithCircuitBreaker(5, 60*time.Second).
    WithMetrics(true)

// Use them anywhere
func GetUser(id int) (*User, error) {
    var user User
    err := userAPI.JSON("GET", fmt.Sprintf("/users/%d", id), nil, &user)
    return &user, err
}
```

## Built-in Features (No Configuration Needed)

✅ **Connection Pooling** - Automatic connection reuse  
✅ **Automatic Retries** - Exponential backoff for failed requests  
✅ **JSON Handling** - Automatic marshaling/unmarshaling  
✅ **Error Handling** - Clear, actionable error messages  
✅ **Timeouts** - Prevents hanging requests  
✅ **HTTP/2 Support** - Automatic protocol negotiation  
✅ **TLS/SSL** - Secure connections by default  

## Optional Advanced Features

🚀 **Rate Limiting** - Configurable requests per second  
🚀 **Circuit Breaker** - Automatic failure detection and recovery  
🚀 **Response Caching** - TTL-based response caching  
🚀 **Load Balancing** - Multiple strategies (round-robin, random, least-conn)  
🚀 **Health Checks** - Automatic endpoint health monitoring  
🚀 **Request Compression** - Automatic gzip compression  
🚀 **Request Signing** - RSA signature support  
🚀 **IP Whitelisting** - Network-level access control  
🚀 **Backup Endpoints** - Automatic failover support  
🚀 **Custom Transport** - Pluggable transport layer  
🚀 **Cookie Management** - Automatic cookie handling  
🚀 **Redirect Control** - Configurable redirect policies  
🚀 **Request/Response Interceptors** - Middleware support  
🚀 **Prometheus Metrics** - Request metrics and monitoring  
🚀 **OpenTelemetry Tracing** - Distributed tracing support  
🚀 **Debug Logging** - Detailed request/response logging  
🧠 **AI-Powered Retry** - Machine learning retry strategies  
🧠 **Smart Caching** - Intelligent caching decisions  
🧠 **Predictive Preloading** - Anticipate future requests  
🧠 **Adaptive Timeouts** - Dynamic timeout optimization  
⚡ **HTTP Streaming** - Real-time data streaming  
⚡ **WebSocket Support** - Full-duplex communication  
⚡ **Server-Sent Events** - Real-time event streaming  
⚡ **Batch Operations** - Concurrent request execution  
⚡ **Pipeline Operations** - Sequential request streaming  
🎯 **GraphQL Support** - Native GraphQL client  
🔧 **Built-in Mocking** - Development and testing support  
🔧 **Record/Replay** - Request recording and playback  
🔧 **Chaos Engineering** - Fault injection testing  
🔧 **Request Validation** - Schema-based validation  
🌐 **HTTP/3 Support** - Next-generation HTTP protocol  
🌐 **Multipath TCP** - Enhanced connection reliability  
🌐 **DNS over HTTPS** - Secure DNS resolution  

## Examples

### Basic API Client

```go
package main

import (
    "fmt"
    "github.com/yourorg/httpclient"
)

type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

func main() {
    // Get a user
    var user User
    err := httpclient.JSON("GET", "https://jsonplaceholder.typicode.com/users/1", nil, &user)
    if err != nil {
        panic(err)
    }
    fmt.Printf("User: %+v\n", user)
}
```

### AI-Powered Client

```go
client := httpclient.New().
    WithAIRetry(true).                    // Smart retry with ML
    WithSmartCaching(true).               // Intelligent caching
    WithPredictivePreloading(true).       // Preload likely requests
    WithAdaptiveTimeout(true).            // Dynamic timeouts
    WithPerformanceOptimization(true)     // Real-time optimization

var user User
err := client.JSON("GET", "https://api.example.com/users/1", nil, &user)
```

### Real-time Features

```go
// WebSocket
ws, err := httpclient.WebSocket("wss://api.example.com/ws")
ws.Send("Hello!")
data, err := ws.Receive()

// Streaming
stream, err := httpclient.Stream("GET", "https://api.example.com/events", nil)
for data := range stream {
    fmt.Printf("Event: %s\n", data)
}

// Batch requests
responses, err := httpclient.Batch().
    Add("GET", "https://api.example.com/users/1", nil).
    Add("GET", "https://api.example.com/users/2", nil).
    Execute()
```
### Production API Client

```go
type APIClient struct {
    client httpclient.Client
}

func NewAPIClient(baseURL, token string) *APIClient {
    client := httpclient.New().
        WithBaseURL(baseURL).
        WithAuth(token).
        WithTimeout(30 * time.Second).
        WithRetries(3).
        WithRateLimiter(50).
        WithCircuitBreaker(5, 60*time.Second).
        WithMetrics(true).
        WithTracing(true)

    return &APIClient{client: client}
}

func (c *APIClient) GetUser(id int) (*User, error) {
    var user User
    err := c.client.JSON("GET", fmt.Sprintf("/users/%d", id), nil, &user)
    return &user, err
}

func (c *APIClient) CreateUser(user *User) (*User, error) {
    var created User
    err := c.client.JSON("POST", "/users", user, &created)
    return &created, err
}
```

## Error Handling

Errors are clear and actionable:

```go
data, err := httpclient.GET("https://api.example.com/users")
if err != nil {
    fmt.Printf("Request failed: %v\n", err)
    // Output: "HTTP 404: User not found"
}
```

## Testing

Easy to test with httptest:

```go
func TestGetUser(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte(`{"id": 1, "name": "John"}`))
    }))
    defer server.Close()
    
    var user User
    err := httpclient.JSON("GET", server.URL, nil, &user)
    assert.NoError(t, err)
    assert.Equal(t, "John", user.Name)
}
```

## Package Structure

```
httpclient/
├── httpclient.go           # Main package interface
├── internal/
│   ├── client/            # Core client implementation
│   ├── config/            # Configuration management
│   ├── middleware/        # Middleware (metrics, tracing, debug)
│   └── retry/             # Retry strategies
├── examples/              # Usage examples
│   ├── basic/            # Basic usage examples
│   ├── advanced/         # Advanced features
│   └── microservice/     # Microservice patterns
└── test/                 # Comprehensive tests
```

## Performance

- **Connection Pooling**: Reuses connections automatically
- **HTTP/2**: Automatic protocol negotiation
- **Keep-Alive**: Persistent connections
- **Compression**: Automatic gzip handling
- **Rate Limiting**: Prevents overwhelming servers
- **Circuit Breaker**: Fails fast when services are down
- **AI Optimization**: Machine learning performance tuning
- **Smart Caching**: Intelligent cache decisions
- **Predictive Loading**: Anticipates future requests
- **Adaptive Timeouts**: Dynamic timeout optimization

## Observability

### Prometheus Metrics

```go
client := httpclient.New().WithMetrics(true)
// Automatically exports:
// - httpclient_requests_total
// - httpclient_request_duration_seconds
```

### OpenTelemetry Tracing

```go
client := httpclient.New().WithTracing(true)
// Automatically creates spans for all requests
```

### Debug Logging

```go
client := httpclient.New().WithDebug(true)
// Logs all requests and responses
```

### Real-time Metrics

```go
client := httpclient.New().WithRealTimeMetrics(true)
// Live performance metrics and dashboards
```
## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Why Choose This Client?

- **🚀 Zero Learning Curve**: If you know HTTP, you know this client
- **🏗️ Production Ready**: Used in production systems handling millions of requests
- **🔧 Highly Configurable**: Every aspect can be customized
- **📊 Observable**: Built-in metrics and tracing
- **🛡️ Reliable**: Automatic retries, circuit breakers, and error handling
- **⚡ Fast**: Optimized for performance with connection pooling
- **🧪 Testable**: Easy to mock and test
- **📚 Well Documented**: Comprehensive examples and documentation
- **🧠 Intelligent**: AI-powered features for optimal performance
- **⚡ Real-time**: WebSocket, streaming, and real-time capabilities
- **🎯 Modern**: GraphQL, HTTP/3, and cutting-edge protocols
- **🔧 Developer-Friendly**: Built-in mocking, testing, and debugging tools

Start simple, scale to enterprise. That's the httpclient way.
