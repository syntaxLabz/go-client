// Package httpclient provides an ultra-simple yet powerful HTTP client for Go.
// It's designed to be as easy as a single function call while providing
// enterprise-grade features under the hood.
//
// Basic usage:
//
//	data, err := httpclient.GET("https://api.example.com/users")
//	err := httpclient.JSON("POST", "https://api.example.com/users", user, &result)
//
// Advanced usage:
//
//	client := httpclient.New().
//		WithAuth("token").
//		WithRetries(5).
//		WithTimeout(30*time.Second)
//
//	data, err := client.GET("https://api.example.com/protected")
package httpclient

import (
	"context"
	"time"

	"github.com/yourorg/httpclient/internal/client"
	"github.com/yourorg/httpclient/internal/config"
)

// Default client instance - ready to use immediately
var Default = New()

// Client is the main HTTP client interface
type Client interface {
	// HTTP Methods
	GET(url string) ([]byte, error)
	POST(url string, body interface{}) ([]byte, error)
	PUT(url string, body interface{}) ([]byte, error)
	PATCH(url string, body interface{}) ([]byte, error)
	DELETE(url string) ([]byte, error)
	HEAD(url string) error
	OPTIONS(url string) ([]byte, error)

	// Context-aware methods
	GetContext(ctx context.Context, url string) ([]byte, error)
	PostContext(ctx context.Context, url string, body interface{}) ([]byte, error)
	PutContext(ctx context.Context, url string, body interface{}) ([]byte, error)
	PatchContext(ctx context.Context, url string, body interface{}) ([]byte, error)
	DeleteContext(ctx context.Context, url string) ([]byte, error)

	// JSON methods
	JSON(method, url string, body, result interface{}) error
	JSONContext(ctx context.Context, method, url string, body, result interface{}) error

	// Streaming methods
	Stream(method, url string, body interface{}) (<-chan []byte, error)
	StreamContext(ctx context.Context, method, url string, body interface{}) (<-chan []byte, error)

	// Batch operations
	Batch() BatchRequest
	Pipeline() PipelineRequest

	// WebSocket support
	WebSocket(url string) (WebSocketConn, error)
	WebSocketContext(ctx context.Context, url string) (WebSocketConn, error)

	// GraphQL support
	GraphQL(query string, variables map[string]interface{}, result interface{}) error
	GraphQLContext(ctx context.Context, query string, variables map[string]interface{}, result interface{}) error

	// Configuration methods (fluent interface)
	WithTimeout(timeout time.Duration) Client
	WithRetries(retries int) Client
	WithBaseURL(baseURL string) Client
	WithAuth(token string) Client
	WithAPIKey(key, value string) Client
	WithHeader(key, value string) Client
	WithHeaders(headers map[string]string) Client
	WithUserAgent(userAgent string) Client
	WithRateLimiter(rps int) Client
	WithCircuitBreaker(threshold int, timeout time.Duration) Client
	WithCache(ttl time.Duration) Client
	WithMetrics(enabled bool) Client
	WithTracing(enabled bool) Client
	WithDebug(enabled bool) Client

	// Advanced features
	WithLoadBalancer(endpoints []string, strategy string) Client
	WithHealthCheck(interval time.Duration, endpoint string) Client
	WithCompression(enabled bool) Client
	WithRequestSigning(keyID, privateKey string) Client
	WithIPWhitelist(ips []string) Client
	WithRequestInterceptor(interceptor func(*http.Request) error) Client
	WithResponseInterceptor(interceptor func(*http.Response) error) Client
	WithBackupEndpoints(endpoints []string) Client
	WithCustomTransport(transport http.RoundTripper) Client
	WithConnectionPool(maxIdle, maxIdlePerHost int) Client
	WithKeepAlive(duration time.Duration) Client
	WithTLSConfig(config *tls.Config) Client
	WithProxy(proxyURL string) Client
	WithCookieJar(jar http.CookieJar) Client
	WithRedirectPolicy(policy func(req *http.Request, via []*http.Request) error) Client

	// AI/ML Features
	WithAIRetry(enabled bool) Client
	WithSmartCaching(enabled bool) Client
	WithPredictivePreloading(enabled bool) Client
	WithAdaptiveTimeout(enabled bool) Client

	// Advanced Networking
	WithHTTP3(enabled bool) Client
	WithMultipath(enabled bool) Client
	WithDNSOverHTTPS(enabled bool) Client
	WithEdgeOptimization(enabled bool) Client

	// Security & Compliance
	WithMTLS(certFile, keyFile string) Client
	WithOAuth2(config OAuth2Config) Client
	WithJWT(config JWTConfig) Client
	WithAPIGateway(config APIGatewayConfig) Client
	WithCompliance(standards []string) Client

	// Performance & Monitoring
	WithRealTimeMetrics(enabled bool) Client
	WithAPM(provider string) Client
	WithChaosEngineering(config ChaosConfig) Client
	WithPerformanceOptimization(enabled bool) Client

	// Developer Experience
	WithMocking(enabled bool) Client
	WithRecording(enabled bool) Client
	WithReplay(enabled bool) Client
	WithValidation(schema interface{}) Client
	WithAutoRetry(config AutoRetryConfig) Client
}

// Advanced types for new features
type BatchRequest interface {
	Add(method, url string, body interface{}) BatchRequest
	Execute() ([]BatchResponse, error)
	ExecuteContext(ctx context.Context) ([]BatchResponse, error)
}

type PipelineRequest interface {
	Add(method, url string, body interface{}) PipelineRequest
	Execute() (<-chan PipelineResponse, error)
	ExecuteContext(ctx context.Context) (<-chan PipelineResponse, error)
}

type BatchResponse struct {
	Index    int
	Data     []byte
	Error    error
	Duration time.Duration
}

type PipelineResponse struct {
	Index    int
	Data     []byte
	Error    error
	Duration time.Duration
}

type WebSocketConn interface {
	Send(data interface{}) error
	Receive() ([]byte, error)
	Close() error
}

type OAuth2Config struct {
	ClientID     string
	ClientSecret string
	TokenURL     string
	Scopes       []string
}

type JWTConfig struct {
	Secret    string
	Algorithm string
	Issuer    string
}

type APIGatewayConfig struct {
	Provider string
	Config   map[string]interface{}
}

type ChaosConfig struct {
	FailureRate    float64
	LatencyRange   [2]time.Duration
	ErrorTypes     []string
	EnabledMethods []string
}

type AutoRetryConfig struct {
	MaxAttempts     int
	BackoffStrategy string
	RetryConditions []string
	JitterEnabled   bool
}
// New creates a new HTTP client with sensible defaults
func New() Client {
	return client.New(config.Default())
}

// NewWithConfig creates a new HTTP client with custom configuration
func NewWithConfig(cfg *config.Config) Client {
	return client.New(cfg)
}

// Smart constructors for common use cases
func NewForMicroservices() Client {
	return New().
		WithLoadBalancer([]string{}, "least-conn").
		WithHealthCheck(30*time.Second, "/health").
		WithCircuitBreaker(5, 60*time.Second).
		WithMetrics(true).
		WithTracing(true).
		WithAIRetry(true).
		WithAdaptiveTimeout(true)
}

func NewForAPI() Client {
	return New().
		WithRateLimiter(100).
		WithCache(5*time.Minute).
		WithCompression(true).
		WithSmartCaching(true).
		WithPerformanceOptimization(true)
}

func NewForEnterprise() Client {
	return New().
		WithLoadBalancer([]string{}, "round-robin").
		WithHealthCheck(30*time.Second, "/health").
		WithCircuitBreaker(10, 120*time.Second).
		WithCompression(true).
		WithMetrics(true).
		WithTracing(true).
		WithRealTimeMetrics(true).
		WithCompliance([]string{"SOC2", "GDPR", "HIPAA"}).
		WithPerformanceOptimization(true).
		WithAIRetry(true).
		WithAdaptiveTimeout(true)
}

func NewForDevelopment() Client {
	return New().
		WithDebug(true).
		WithMocking(true).
		WithRecording(true).
		WithValidation(nil).
		WithChaosEngineering(ChaosConfig{
			FailureRate: 0.1,
			LatencyRange: [2]time.Duration{100*time.Millisecond, 500*time.Millisecond},
		})
}
// Package-level convenience functions using the default client

// GET makes a GET request using the default client
func GET(url string) ([]byte, error) {
	return Default.GET(url)
}

// POST makes a POST request using the default client
func POST(url string, body interface{}) ([]byte, error) {
	return Default.POST(url, body)
}

// PUT makes a PUT request using the default client
func PUT(url string, body interface{}) ([]byte, error) {
	return Default.PUT(url, body)
}

// PATCH makes a PATCH request using the default client
func PATCH(url string, body interface{}) ([]byte, error) {
	return Default.PATCH(url, body)
}

// DELETE makes a DELETE request using the default client
func DELETE(url string) ([]byte, error) {
	return Default.DELETE(url)
}

// HEAD makes a HEAD request using the default client
func HEAD(url string) error {
	return Default.HEAD(url)
}

// OPTIONS makes an OPTIONS request using the default client
func OPTIONS(url string) ([]byte, error) {
	return Default.OPTIONS(url)
}

// JSON makes a JSON request using the default client
func JSON(method, url string, body, result interface{}) error {
	return Default.JSON(method, url, body, result)
}

// Advanced package-level functions
func Batch() BatchRequest {
	return Default.Batch()
}

func Pipeline() PipelineRequest {
	return Default.Pipeline()
}

func Stream(method, url string, body interface{}) (<-chan []byte, error) {
	return Default.Stream(method, url, body)
}

func GraphQL(query string, variables map[string]interface{}, result interface{}) error {
	return Default.GraphQL(query, variables, result)
}

func WebSocket(url string) (WebSocketConn, error) {
	return Default.WebSocket(url)
}
// Context-aware package-level functions

// GetContext makes a GET request with context using the default client
func GetContext(ctx context.Context, url string) ([]byte, error) {
	return Default.GetContext(ctx, url)
}

// PostContext makes a POST request with context using the default client
func PostContext(ctx context.Context, url string, body interface{}) ([]byte, error) {
	return Default.PostContext(ctx, url, body)
}

// JSONContext makes a JSON request with context using the default client
func JSONContext(ctx context.Context, method, url string, body, result interface{}) error {
	return Default.JSONContext(ctx, method, url, body, result)
}