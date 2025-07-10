package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

// CircuitState represents the state of the circuit breaker
type CircuitState int

const (
	StateClosed CircuitState = iota
	StateOpen
	StateHalfOpen
)

// CircuitBreaker middleware
type circuitBreakerMiddleware struct {
	state         CircuitState
	failures      int64
	lastFailTime  time.Time
	threshold     int64
	timeout       time.Duration
	mu            sync.RWMutex
}

// NewCircuitBreaker creates a new circuit breaker middleware
func NewCircuitBreaker(threshold int, timeout time.Duration) Middleware {
	return &circuitBreakerMiddleware{
		state:     StateClosed,
		threshold: int64(threshold),
		timeout:   timeout,
	}
}

func (cb *circuitBreakerMiddleware) Before(req *http.Request) error {
	cb.mu.RLock()
	state := cb.state
	failures := cb.failures
	lastFailTime := cb.lastFailTime
	cb.mu.RUnlock()
	
	switch state {
	case StateOpen:
		if time.Since(lastFailTime) > cb.timeout {
			cb.mu.Lock()
			cb.state = StateHalfOpen
			cb.mu.Unlock()
		} else {
			return fmt.Errorf("circuit breaker is open")
		}
	case StateHalfOpen:
		// Allow one request through
	case StateClosed:
		// Normal operation
	}
	
	return nil
}

func (cb *circuitBreakerMiddleware) After(resp *http.Response) {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	
	if resp.StatusCode >= 500 {
		// Server error - count as failure
		cb.failures++
		cb.lastFailTime = time.Now()
		
		if cb.failures >= cb.threshold {
			cb.state = StateOpen
		}
	} else {
		// Success - reset failures
		cb.failures = 0
		if cb.state == StateHalfOpen {
			cb.state = StateClosed
		}
	}
}

// GetState returns the current circuit breaker state
func (cb *circuitBreakerMiddleware) GetState() CircuitState {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// GetFailures returns the current failure count
func (cb *circuitBreakerMiddleware) GetFailures() int64 {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.failures
}