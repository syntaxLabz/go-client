package config

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"time"
)

// Config holds all client configuration options
type Config struct {
	// Basic settings
	Timeout     time.Duration
	BaseURL     string
	UserAgent   string
	Headers     map[string]string

	// Retry settings
	Retries         int
	RetryDelay      time.Duration
	RetryMultiplier float64
	RetryMaxDelay   time.Duration

	// Connection settings
	MaxIdleConns        int
	MaxIdleConnsPerHost int
	IdleConnTimeout     time.Duration
	KeepAlive           time.Duration

	// Rate limiting
	RateLimitRPS int

	// Circuit breaker
	CircuitBreakerThreshold int
	CircuitBreakerTimeout   time.Duration

	// Caching
	CacheTTL time.Duration

	// Observability
	MetricsEnabled bool
	TracingEnabled bool
	DebugEnabled   bool

	// Security
	TLSInsecureSkipVerify bool
	TLSTimeout            time.Duration

	// Advanced features
	LoadBalancerEndpoints []string
	LoadBalancerStrategy  string
	HealthCheckInterval   time.Duration
	HealthCheckEndpoint   string
	CompressionEnabled    bool
	RequestSigningKeyID   string
	RequestSigningKey     string
	IPWhitelist          []string
	BackupEndpoints      []string
	CustomTransport      http.RoundTripper
	TLSConfig            *tls.Config
	ProxyURL             *url.URL
	CookieJar            http.CookieJar
	RedirectPolicy       func(req *http.Request, via []*http.Request) error
	RequestInterceptors  []func(*http.Request) error
	ResponseInterceptors []func(*http.Response) error

	// AI/ML Features
	AIRetryEnabled            bool
	SmartCachingEnabled       bool
	PredictivePreloadingEnabled bool
	AdaptiveTimeoutEnabled    bool

	// Advanced Networking
	HTTP3Enabled           bool
	MultipathEnabled       bool
	DNSOverHTTPSEnabled    bool
	EdgeOptimizationEnabled bool

	// Security & Compliance
	MTLSCertFile        string
	MTLSKeyFile         string
	OAuth2Config        *OAuth2Config
	JWTConfig           *JWTConfig
	APIGatewayConfig    *APIGatewayConfig
	ComplianceStandards []string

	// Performance & Monitoring
	RealTimeMetricsEnabled      bool
	APMProvider                 string
	ChaosEngineeringEnabled     bool
	ChaosConfig                 *ChaosConfig
	PerformanceOptimizationEnabled bool

	// Developer Experience
	MockingEnabled    bool
	RecordingEnabled  bool
	ReplayEnabled     bool
	ValidationSchema  interface{}
	AutoRetryConfig   *AutoRetryConfig

	// Streaming & Real-time
	StreamingEnabled    bool
	WebSocketEnabled    bool
	ServerSentEventsEnabled bool

	// GraphQL
	GraphQLEnabled bool
	GraphQLEndpoint string

	// Batch & Pipeline
	BatchEnabled    bool
	PipelineEnabled bool
	MaxBatchSize    int
	MaxPipelineSize int
}

// Advanced configuration types
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
}

// Default returns a configuration with sensible defaults
func Default() *Config {
	return &Config{
		// Basic settings
		Timeout:   30 * time.Second,
		UserAgent: "httpclient/1.0",
		Headers:   make(map[string]string),

		// Retry settings
		Retries:         3,
		RetryDelay:      1 * time.Second,
		RetryMultiplier: 2.0,
		RetryMaxDelay:   30 * time.Second,

		// Connection settings
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
		KeepAlive:           30 * time.Second,

		// Rate limiting
		RateLimitRPS: 100,

		// Circuit breaker
		CircuitBreakerThreshold: 5,
		CircuitBreakerTimeout:   60 * time.Second,

		// Caching
		CacheTTL: 5 * time.Minute,

		// Observability
		MetricsEnabled: false,
		TracingEnabled: false,
		DebugEnabled:   false,

		// Security
		TLSInsecureSkipVerify: false,
		TLSTimeout:            10 * time.Second,

		// AI/ML Features (enabled by default for smart behavior)
		AIRetryEnabled:              true,
		SmartCachingEnabled:         true,
		PredictivePreloadingEnabled: false,
		AdaptiveTimeoutEnabled:      true,

		// Advanced Networking
		HTTP3Enabled:           false,
		MultipathEnabled:       false,
		DNSOverHTTPSEnabled:    false,
		EdgeOptimizationEnabled: false,

		// Performance & Monitoring
		RealTimeMetricsEnabled:         false,
		PerformanceOptimizationEnabled: true,

		// Developer Experience
		MockingEnabled:   false,
		RecordingEnabled: false,
		ReplayEnabled:    false,

		// Streaming & Real-time
		StreamingEnabled:        true,
		WebSocketEnabled:        true,
		ServerSentEventsEnabled: true,

		// GraphQL
		GraphQLEnabled: true,

		// Batch & Pipeline
		BatchEnabled:    true,
		PipelineEnabled: true,
		MaxBatchSize:    100,
		MaxPipelineSize: 50,
	}
}

// Clone creates a deep copy of the configuration
func (c *Config) Clone() *Config {
	clone := *c
	clone.Headers = make(map[string]string)
	for k, v := range c.Headers {
		clone.Headers[k] = v
	}

	// Clone complex types
	if c.OAuth2Config != nil {
		oauth2Clone := *c.OAuth2Config
		clone.OAuth2Config = &oauth2Clone
	}
	if c.JWTConfig != nil {
		jwtClone := *c.JWTConfig
		clone.JWTConfig = &jwtClone
	}
	if c.APIGatewayConfig != nil {
		gatewayClone := *c.APIGatewayConfig
		clone.APIGatewayConfig = &gatewayClone
	}
	if c.ChaosConfig != nil {
		chaosClone := *c.ChaosConfig
		clone.ChaosConfig = &chaosClone
	}
	if c.AutoRetryConfig != nil {
		retryClone := *c.AutoRetryConfig
		clone.AutoRetryConfig = &retryClone
	}

	return &clone
}