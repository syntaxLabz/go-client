package retry

import (
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/yourorg/httpclient/internal/config"
)

// Strategy defines the retry strategy interface
type Strategy interface {
	Execute(fn func() ([]byte, error)) ([]byte, error)
}

// exponentialBackoff implements exponential backoff retry strategy
type exponentialBackoff struct {
	maxRetries  int
	baseDelay   time.Duration
	multiplier  float64
	maxDelay    time.Duration
}

// NewExponentialBackoff creates a new exponential backoff retry strategy
func NewExponentialBackoff(cfg *config.Config) Strategy {
	return &exponentialBackoff{
		maxRetries: cfg.Retries,
		baseDelay:  cfg.RetryDelay,
		multiplier: cfg.RetryMultiplier,
		maxDelay:   cfg.RetryMaxDelay,
	}
}

func (e *exponentialBackoff) Execute(fn func() ([]byte, error)) ([]byte, error) {
	var lastErr error
	
	for attempt := 0; attempt <= e.maxRetries; attempt++ {
		data, err := fn()
		if err == nil {
			return data, nil
		}
		
		lastErr = err
		
		// Don't retry on client errors (4xx)
		if httpErr, ok := err.(*HTTPError); ok {
			if httpErr.StatusCode >= 400 && httpErr.StatusCode < 500 {
				return nil, err
			}
		}
		
		// Don't sleep after the last attempt
		if attempt < e.maxRetries {
			delay := e.calculateDelay(attempt)
			time.Sleep(delay)
		}
	}
	
	return nil, fmt.Errorf("max retries exceeded: %w", lastErr)
}

func (e *exponentialBackoff) calculateDelay(attempt int) time.Duration {
	delay := float64(e.baseDelay) * math.Pow(e.multiplier, float64(attempt))
	if delay > float64(e.maxDelay) {
		delay = float64(e.maxDelay)
	}
	return time.Duration(delay)
}

// HTTPError represents an HTTP error with status code
type HTTPError struct {
	StatusCode int
	Message    string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTP %d: %s", e.StatusCode, e.Message)
}

// NewHTTPError creates a new HTTP error
func NewHTTPError(statusCode int, message string) *HTTPError {
	return &HTTPError{
		StatusCode: statusCode,
		Message:    message,
	}
}