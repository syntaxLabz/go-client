package middleware

import (
	"net/http"
)

// Middleware defines the interface for HTTP middleware
type Middleware interface {
	Before(req *http.Request) error
	After(resp *http.Response)
}

// middlewareFunc is an adapter to allow functions to be used as middleware
type middlewareFunc struct {
	before func(*http.Request) error
	after  func(*http.Response)
}

func (m middlewareFunc) Before(req *http.Request) error {
	if m.before != nil {
		return m.before(req)
	}
	return nil
}

func (m middlewareFunc) After(resp *http.Response) {
	if m.after != nil {
		m.after(resp)
	}
}

// NewFunc creates a middleware from functions
func NewFunc(before func(*http.Request) error, after func(*http.Response)) Middleware {
	return middlewareFunc{before: before, after: after}
}