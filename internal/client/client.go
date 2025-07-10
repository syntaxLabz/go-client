package client

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/yourorg/httpclient/internal/config"
	"github.com/yourorg/httpclient/internal/loadbalancer"
	"github.com/yourorg/httpclient/internal/middleware"
	"github.com/yourorg/httpclient/internal/retry"
	"golang.org/x/time/rate"
)

// client implements the Client interface
type client struct {
	httpClient     *http.Client
	config         *config.Config
	rateLimiter    *rate.Limiter
	middlewares    []middleware.Middleware
	retryStrategy  retry.Strategy
	loadBalancer   loadbalancer.LoadBalancer
	healthChecker  *HealthChecker
	requestSigner  *RequestSigner
	ipWhitelist    map[string]bool
	backupClients  []*client
	mu             sync.RWMutex
}

// HealthChecker manages endpoint health checking
type HealthChecker struct {
	endpoints map[string]*EndpointHealth
	interval  time.Duration
	client    *http.Client
	mu        sync.RWMutex
}

type EndpointHealth struct {
	URL       string
	Healthy   bool
	LastCheck time.Time
	Failures  int64
}

// RequestSigner handles request signing
type RequestSigner struct {
	keyID      string
	privateKey *rsa.PrivateKey
}

// New creates a new HTTP client with the given configuration
func New(cfg *config.Config) *client {
	var transport http.RoundTripper
	
	if cfg.CustomTransport != nil {
		transport = cfg.CustomTransport
	} else {
		tlsConfig := cfg.TLSConfig
		if tlsConfig == nil {
			tlsConfig = &tls.Config{
				InsecureSkipVerify: cfg.TLSInsecureSkipVerify,
			}
		}

		httpTransport := &http.Transport{
			MaxIdleConns:        cfg.MaxIdleConns,
			MaxIdleConnsPerHost: cfg.MaxIdleConnsPerHost,
			IdleConnTimeout:     cfg.IdleConnTimeout,
			TLSClientConfig:     tlsConfig,
			TLSHandshakeTimeout: cfg.TLSTimeout,
			KeepAlive:          cfg.KeepAlive,
		}

		if cfg.ProxyURL != nil {
			httpTransport.Proxy = http.ProxyURL(cfg.ProxyURL)
		}

		if cfg.CompressionEnabled {
			transport = &compressionTransport{base: httpTransport}
		} else {
			transport = httpTransport
		}
	}

	httpClient := &http.Client{
		Timeout:       cfg.Timeout,
		Transport:     transport,
		Jar:          cfg.CookieJar,
		CheckRedirect: cfg.RedirectPolicy,
	}

	var rateLimiter *rate.Limiter
	if cfg.RateLimitRPS > 0 {
		rateLimiter = rate.NewLimiter(rate.Limit(cfg.RateLimitRPS), cfg.RateLimitRPS)
	}

	// Initialize load balancer
	var lb loadbalancer.LoadBalancer
	if len(cfg.LoadBalancerEndpoints) > 0 {
		lb = loadbalancer.New(cfg.LoadBalancerEndpoints, cfg.LoadBalancerStrategy)
	}

	// Initialize health checker
	var hc *HealthChecker
	if cfg.HealthCheckInterval > 0 && cfg.HealthCheckEndpoint != "" {
		hc = NewHealthChecker(cfg.HealthCheckInterval, cfg.HealthCheckEndpoint)
		go hc.Start()
	}

	// Initialize request signer
	var rs *RequestSigner
	if cfg.RequestSigningKeyID != "" && cfg.RequestSigningKey != "" {
		if signer, err := NewRequestSigner(cfg.RequestSigningKeyID, cfg.RequestSigningKey); err == nil {
			rs = signer
		}
	}

	// Initialize IP whitelist
	ipWhitelist := make(map[string]bool)
	for _, ip := range cfg.IPWhitelist {
		ipWhitelist[ip] = true
	}

	c := &client{
		httpClient:     httpClient,
		config:         cfg,
		rateLimiter:    rateLimiter,
		middlewares:    []middleware.Middleware{},
		retryStrategy:  retry.NewExponentialBackoff(cfg),
		loadBalancer:   lb,
		healthChecker:  hc,
		requestSigner:  rs,
		ipWhitelist:    ipWhitelist,
	}

	// Initialize backup clients
	for _, endpoint := range cfg.BackupEndpoints {
		backupCfg := cfg.Clone()
		backupCfg.BaseURL = endpoint
		c.backupClients = append(c.backupClients, New(backupCfg))
	}

	// Add default middlewares
	if cfg.MetricsEnabled {
		c.middlewares = append(c.middlewares, middleware.NewMetrics())
	}
	if cfg.TracingEnabled {
		c.middlewares = append(c.middlewares, middleware.NewTracing())
	}
	if cfg.DebugEnabled {
		c.middlewares = append(c.middlewares, middleware.NewDebug())
	}

	return c
}

// HTTP Methods

func (c *client) GET(url string) ([]byte, error) {
	return c.GetContext(context.Background(), url)
}

func (c *client) POST(url string, body interface{}) ([]byte, error) {
	return c.PostContext(context.Background(), url, body)
}

func (c *client) PUT(url string, body interface{}) ([]byte, error) {
	return c.PutContext(context.Background(), url, body)
}

func (c *client) PATCH(url string, body interface{}) ([]byte, error) {
	return c.PatchContext(context.Background(), url, body)
}

func (c *client) DELETE(url string) ([]byte, error) {
	return c.DeleteContext(context.Background(), url)
}

func (c *client) HEAD(url string) error {
	_, err := c.do(context.Background(), "HEAD", url, nil)
	return err
}

func (c *client) OPTIONS(url string) ([]byte, error) {
	return c.do(context.Background(), "OPTIONS", url, nil)
}

// Context-aware methods

func (c *client) GetContext(ctx context.Context, url string) ([]byte, error) {
	return c.do(ctx, "GET", url, nil)
}

func (c *client) PostContext(ctx context.Context, url string, body interface{}) ([]byte, error) {
	return c.do(ctx, "POST", url, body)
}

func (c *client) PutContext(ctx context.Context, url string, body interface{}) ([]byte, error) {
	return c.do(ctx, "PUT", url, body)
}

func (c *client) PatchContext(ctx context.Context, url string, body interface{}) ([]byte, error) {
	return c.do(ctx, "PATCH", url, body)
}

func (c *client) DeleteContext(ctx context.Context, url string) ([]byte, error) {
	return c.do(ctx, "DELETE", url, nil)
}

// JSON methods

func (c *client) JSON(method, url string, body, result interface{}) error {
	return c.JSONContext(context.Background(), method, url, body, result)
}

func (c *client) JSONContext(ctx context.Context, method, url string, body, result interface{}) error {
	data, err := c.do(ctx, method, url, body)
	if err != nil {
		return err
	}
	if result != nil && len(data) > 0 {
		return json.Unmarshal(data, result)
	}
	return nil
}

// Configuration methods (fluent interface)

func (c *client) WithTimeout(timeout time.Duration) *client {
	newConfig := c.config.Clone()
	newConfig.Timeout = timeout
	return New(newConfig)
}

func (c *client) WithRetries(retries int) *client {
	newConfig := c.config.Clone()
	newConfig.Retries = retries
	return New(newConfig)
}

func (c *client) WithBaseURL(baseURL string) *client {
	newConfig := c.config.Clone()
	newConfig.BaseURL = strings.TrimSuffix(baseURL, "/")
	return New(newConfig)
}

func (c *client) WithAuth(token string) *client {
	return c.WithHeader("Authorization", "Bearer "+token)
}

func (c *client) WithAPIKey(key, value string) *client {
	return c.WithHeader(key, value)
}

func (c *client) WithHeader(key, value string) *client {
	newConfig := c.config.Clone()
	newConfig.Headers[key] = value
	return New(newConfig)
}

func (c *client) WithHeaders(headers map[string]string) *client {
	newConfig := c.config.Clone()
	for k, v := range headers {
		newConfig.Headers[k] = v
	}
	return New(newConfig)
}

func (c *client) WithUserAgent(userAgent string) *client {
	newConfig := c.config.Clone()
	newConfig.UserAgent = userAgent
	return New(newConfig)
}

func (c *client) WithRateLimiter(rps int) *client {
	newConfig := c.config.Clone()
	newConfig.RateLimitRPS = rps
	return New(newConfig)
}

func (c *client) WithCircuitBreaker(threshold int, timeout time.Duration) *client {
	newConfig := c.config.Clone()
	newConfig.CircuitBreakerThreshold = threshold
	newConfig.CircuitBreakerTimeout = timeout
	return New(newConfig)
}

func (c *client) WithCache(ttl time.Duration) *client {
	newConfig := c.config.Clone()
	newConfig.CacheTTL = ttl
	return New(newConfig)
}

func (c *client) WithMetrics(enabled bool) *client {
	newConfig := c.config.Clone()
	newConfig.MetricsEnabled = enabled
	return New(newConfig)
}

func (c *client) WithTracing(enabled bool) *client {
	newConfig := c.config.Clone()
	newConfig.TracingEnabled = enabled
	return New(newConfig)
}

func (c *client) WithDebug(enabled bool) *client {
	newConfig := c.config.Clone()
	newConfig.DebugEnabled = enabled
	return New(newConfig)
}

// Advanced configuration methods

func (c *client) WithLoadBalancer(endpoints []string, strategy string) *client {
	newConfig := c.config.Clone()
	newConfig.LoadBalancerEndpoints = endpoints
	newConfig.LoadBalancerStrategy = strategy
	return New(newConfig)
}

func (c *client) WithHealthCheck(interval time.Duration, endpoint string) *client {
	newConfig := c.config.Clone()
	newConfig.HealthCheckInterval = interval
	newConfig.HealthCheckEndpoint = endpoint
	return New(newConfig)
}

func (c *client) WithCompression(enabled bool) *client {
	newConfig := c.config.Clone()
	newConfig.CompressionEnabled = enabled
	return New(newConfig)
}

func (c *client) WithRequestSigning(keyID, privateKey string) *client {
	newConfig := c.config.Clone()
	newConfig.RequestSigningKeyID = keyID
	newConfig.RequestSigningKey = privateKey
	return New(newConfig)
}

func (c *client) WithIPWhitelist(ips []string) *client {
	newConfig := c.config.Clone()
	newConfig.IPWhitelist = ips
	return New(newConfig)
}

func (c *client) WithRequestInterceptor(interceptor func(*http.Request) error) *client {
	newConfig := c.config.Clone()
	newConfig.RequestInterceptors = append(newConfig.RequestInterceptors, interceptor)
	return New(newConfig)
}

func (c *client) WithResponseInterceptor(interceptor func(*http.Response) error) *client {
	newConfig := c.config.Clone()
	newConfig.ResponseInterceptors = append(newConfig.ResponseInterceptors, interceptor)
	return New(newConfig)
}

func (c *client) WithBackupEndpoints(endpoints []string) *client {
	newConfig := c.config.Clone()
	newConfig.BackupEndpoints = endpoints
	return New(newConfig)
}

func (c *client) WithCustomTransport(transport http.RoundTripper) *client {
	newConfig := c.config.Clone()
	newConfig.CustomTransport = transport
	return New(newConfig)
}

func (c *client) WithConnectionPool(maxIdle, maxIdlePerHost int) *client {
	newConfig := c.config.Clone()
	newConfig.MaxIdleConns = maxIdle
	newConfig.MaxIdleConnsPerHost = maxIdlePerHost
	return New(newConfig)
}

func (c *client) WithKeepAlive(duration time.Duration) *client {
	newConfig := c.config.Clone()
	newConfig.KeepAlive = duration
	return New(newConfig)
}

func (c *client) WithTLSConfig(config *tls.Config) *client {
	newConfig := c.config.Clone()
	newConfig.TLSConfig = config
	return New(newConfig)
}

func (c *client) WithProxy(proxyURL string) *client {
	newConfig := c.config.Clone()
	if u, err := url.Parse(proxyURL); err == nil {
		newConfig.ProxyURL = u
	}
	return New(newConfig)
}

func (c *client) WithCookieJar(jar http.CookieJar) *client {
	newConfig := c.config.Clone()
	newConfig.CookieJar = jar
	return New(newConfig)
}

func (c *client) WithRedirectPolicy(policy func(req *http.Request, via []*http.Request) error) *client {
	newConfig := c.config.Clone()
	newConfig.RedirectPolicy = policy
	return New(newConfig)
}

// Internal methods

func (c *client) do(ctx context.Context, method, urlStr string, body interface{}) ([]byte, error) {
	// Check IP whitelist
	if len(c.ipWhitelist) > 0 {
		if err := c.checkIPWhitelist(urlStr); err != nil {
			return nil, err
		}
	}

	// Rate limiting
	if c.rateLimiter != nil {
		if err := c.rateLimiter.Wait(ctx); err != nil {
			return nil, fmt.Errorf("rate limit exceeded: %w", err)
		}
	}

	// Build URL with load balancing
	fullURL, err := c.buildURLWithLoadBalancing(urlStr)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	// Prepare request body
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(jsonData)
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, method, fullURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	// Set headers
	c.setHeaders(req, body != nil)

	// Apply request interceptors
	for _, interceptor := range c.config.RequestInterceptors {
		if err := interceptor(req); err != nil {
			return nil, fmt.Errorf("request interceptor failed: %w", err)
		}
	}

	// Sign request if configured
	if c.requestSigner != nil {
		if err := c.requestSigner.SignRequest(req); err != nil {
			return nil, fmt.Errorf("request signing failed: %w", err)
		}
	}

	// Execute with retry
	data, err := c.retryStrategy.Execute(func() ([]byte, error) {
		return c.executeRequest(req)
	})

	// Try backup endpoints if primary fails
	if err != nil && len(c.backupClients) > 0 {
		for _, backup := range c.backupClients {
			if backupData, backupErr := backup.do(ctx, method, urlStr, body); backupErr == nil {
				return backupData, nil
			}
		}
	}

	return data, err
}

func (c *client) checkIPWhitelist(urlStr string) error {
	u, err := url.Parse(urlStr)
	if err != nil {
		return err
	}

	host := u.Hostname()
	ips, err := net.LookupIP(host)
	if err != nil {
		return fmt.Errorf("failed to resolve host %s: %w", host, err)
	}

	for _, ip := range ips {
		if c.ipWhitelist[ip.String()] {
			return nil
		}
	}

	return fmt.Errorf("IP not whitelisted for host %s", host)
}

func (c *client) buildURLWithLoadBalancing(urlStr string) (string, error) {
	// Use load balancer if configured
	if c.loadBalancer != nil {
		endpoint := c.loadBalancer.NextEndpoint()
		if endpoint != "" {
			base, err := url.Parse(endpoint)
			if err != nil {
				return "", err
			}
			rel, err := url.Parse(urlStr)
			if err != nil {
				return "", err
			}
			return base.ResolveReference(rel).String(), nil
		}
	}

	// Fallback to base URL
	if c.config.BaseURL == "" {
		return urlStr, nil
	}

	base, err := url.Parse(c.config.BaseURL)
	if err != nil {
		return "", err
	}

	rel, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}

	return base.ResolveReference(rel).String(), nil
}

func (c *client) setHeaders(req *http.Request, hasBody bool) {
	// Set default headers
	req.Header.Set("User-Agent", c.config.UserAgent)
	
	if c.config.CompressionEnabled {
		req.Header.Set("Accept-Encoding", "gzip, deflate")
	}

	if hasBody {
		req.Header.Set("Content-Type", "application/json")
	}

	// Set custom headers
	for key, value := range c.config.Headers {
		req.Header.Set(key, value)
	}
}

func (c *client) executeRequest(req *http.Request) ([]byte, error) {
	// Apply middlewares
	for _, mw := range c.middlewares {
		if err := mw.Before(req); err != nil {
			return nil, err
		}
	}

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Apply middlewares
	for _, mw := range c.middlewares {
		mw.After(resp)
	}

	// Handle compressed responses
	if resp.Header.Get("Content-Encoding") == "gzip" {
		gzipReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("gzip decompression failed: %w", err)
		}
		defer gzipReader.Close()
		resp.Body = gzipReader
	}

	// Read response
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	// Apply response interceptors
	for _, interceptor := range c.config.ResponseInterceptors {
		if err := interceptor(resp); err != nil {
			return nil, fmt.Errorf("response interceptor failed: %w", err)
		}
	}

	// Check status code
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(data))
	}

	return data, nil
}

// Compression transport wrapper
type compressionTransport struct {
	base http.RoundTripper
}

func (ct *compressionTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Body != nil && req.Header.Get("Content-Encoding") == "" {
		// Compress request body
		var buf bytes.Buffer
		gzipWriter := gzip.NewWriter(&buf)
		
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		req.Body.Close()
		
		if _, err := gzipWriter.Write(bodyBytes); err != nil {
			return nil, err
		}
		if err := gzipWriter.Close(); err != nil {
			return nil, err
		}
		
		req.Body = io.NopCloser(&buf)
		req.Header.Set("Content-Encoding", "gzip")
		req.ContentLength = int64(buf.Len())
	}
	
	return ct.base.RoundTrip(req)
}

// Health checker implementation
func NewHealthChecker(interval time.Duration, endpoint string) *HealthChecker {
	return &HealthChecker{
		endpoints: make(map[string]*EndpointHealth),
		interval:  interval,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (hc *HealthChecker) Start() {
	ticker := time.NewTicker(hc.interval)
	defer ticker.Stop()

	for range ticker.C {
		hc.checkEndpoints()
	}
}

func (hc *HealthChecker) checkEndpoints() {
	hc.mu.RLock()
	endpoints := make([]*EndpointHealth, 0, len(hc.endpoints))
	for _, ep := range hc.endpoints {
		endpoints = append(endpoints, ep)
	}
	hc.mu.RUnlock()

	for _, ep := range endpoints {
		go hc.checkEndpoint(ep)
	}
}

func (hc *HealthChecker) checkEndpoint(ep *EndpointHealth) {
	resp, err := hc.client.Get(ep.URL)
	
	hc.mu.Lock()
	defer hc.mu.Unlock()
	
	ep.LastCheck = time.Now()
	
	if err != nil || resp.StatusCode >= 400 {
		ep.Healthy = false
		atomic.AddInt64(&ep.Failures, 1)
	} else {
		ep.Healthy = true
		atomic.StoreInt64(&ep.Failures, 0)
	}
	
	if resp != nil {
		resp.Body.Close()
	}
}

// Request signer implementation
func NewRequestSigner(keyID, privateKeyPEM string) (*RequestSigner, error) {
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return &RequestSigner{
		keyID:      keyID,
		privateKey: privateKey,
	}, nil
}

func (rs *RequestSigner) SignRequest(req *http.Request) error {
	// Create signature string
	sigString := rs.createSignatureString(req)
	
	// Sign the string
	hash := sha256.Sum256([]byte(sigString))
	signature, err := rsa.SignPKCS1v15(rand.Reader, rs.privateKey, crypto.SHA256, hash[:])
	if err != nil {
		return fmt.Errorf("failed to sign request: %w", err)
	}

	// Add signature header
	sigHeader := fmt.Sprintf("keyId=\"%s\",algorithm=\"rsa-sha256\",signature=\"%s\"",
		rs.keyID, base64.StdEncoding.EncodeToString(signature))
	req.Header.Set("Signature", sigHeader)

	return nil
}

func (rs *RequestSigner) createSignatureString(req *http.Request) string {
	var parts []string
	
	// Add method and path
	parts = append(parts, fmt.Sprintf("(request-target): %s %s", 
		strings.ToLower(req.Method), req.URL.RequestURI()))
	
	// Add headers in alphabetical order
	var headerNames []string
	for name := range req.Header {
		headerNames = append(headerNames, strings.ToLower(name))
	}
	sort.Strings(headerNames)
	
	for _, name := range headerNames {
		if name == "signature" {
			continue
		}
		parts = append(parts, fmt.Sprintf("%s: %s", name, req.Header.Get(name)))
	}
	
	return strings.Join(parts, "\n")
}