package middleware

import (
	"fmt"
	"net/http"
	"time"
)

type debugMiddleware struct {
	startTime time.Time
}

// NewDebug creates a new debug middleware
func NewDebug() Middleware {
	return &debugMiddleware{}
}

func (d *debugMiddleware) Before(req *http.Request) error {
	d.startTime = time.Now()
	fmt.Printf("[DEBUG] %s %s\n", req.Method, req.URL.String())
	
	// Print headers
	for key, values := range req.Header {
		for _, value := range values {
			fmt.Printf("[DEBUG] Header: %s: %s\n", key, value)
		}
	}
	
	return nil
}

func (d *debugMiddleware) After(resp *http.Response) {
	duration := time.Since(d.startTime)
	fmt.Printf("[DEBUG] Response: %d %s (took %v)\n", 
		resp.StatusCode, resp.Status, duration)
}