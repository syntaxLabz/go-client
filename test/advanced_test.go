package test

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourorg/httpclient"
)

func TestLoadBalancing(t *testing.T) {
	// Create multiple test servers
	server1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("server1"))
	}))
	defer server1.Close()

	server2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("server2"))
	}))
	defer server2.Close()

	client := httpclient.New().
		WithLoadBalancer([]string{server1.URL, server2.URL}, "round-robin")

	// Make multiple requests to test load balancing
	responses := make(map[string]int)
	for i := 0; i < 10; i++ {
		data, err := client.GET("/")
		if err != nil {
			t.Fatalf("Request %d failed: %v", i, err)
		}
		responses[string(data)]++
	}

	// Both servers should have received requests
	if responses["server1"] == 0 || responses["server2"] == 0 {
		t.Error("Load balancing not working properly")
	}
}

func TestCompression(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Accept-Encoding") == "" {
			t.Error("Expected Accept-Encoding header for compression")
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "compressed response"}`))
	}))
	defer server.Close()

	client := httpclient.New().WithCompression(true)
	data, err := client.GET(server.URL)
	if err != nil {
		t.Fatalf("Compression request failed: %v", err)
	}

	if string(data) != `{"message": "compressed response"}` {
		t.Errorf("Unexpected response: %s", data)
	}
}

func TestRequestInterceptor(t *testing.T) {
	interceptorCalled := false
	
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Test-Header") != "test-value" {
			t.Error("Request interceptor header not found")
		}
		w.Write([]byte("ok"))
	}))
	defer server.Close()

	client := httpclient.New().
		WithRequestInterceptor(func(req *http.Request) error {
			interceptorCalled = true
			req.Header.Set("X-Test-Header", "test-value")
			return nil
		})

	_, err := client.GET(server.URL)
	if err != nil {
		t.Fatalf("Request with interceptor failed: %v", err)
	}

	if !interceptorCalled {
		t.Error("Request interceptor was not called")
	}
}

func TestResponseInterceptor(t *testing.T) {
	interceptorCalled := false
	
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Response-Header", "response-value")
		w.Write([]byte("ok"))
	}))
	defer server.Close()

	client := httpclient.New().
		WithResponseInterceptor(func(resp *http.Response) error {
			interceptorCalled = true
			if resp.Header.Get("X-Response-Header") != "response-value" {
				t.Error("Expected response header not found")
			}
			return nil
		})

	_, err := client.GET(server.URL)
	if err != nil {
		t.Fatalf("Request with response interceptor failed: %v", err)
	}

	if !interceptorCalled {
		t.Error("Response interceptor was not called")
	}
}

func TestBackupEndpoints(t *testing.T) {
	// Primary server that fails
	primaryServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("primary failed"))
	}))
	defer primaryServer.Close()

	// Backup server that works
	backupServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("backup success"))
	}))
	defer backupServer.Close()

	client := httpclient.New().
		WithBaseURL(primaryServer.URL).
		WithBackupEndpoints([]string{backupServer.URL}).
		WithRetries(1)

	data, err := client.GET("/")
	if err != nil {
		t.Fatalf("Backup endpoint request failed: %v", err)
	}

	if string(data) != "backup success" {
		t.Errorf("Expected backup response, got: %s", data)
	}
}

func TestTLSConfig(t *testing.T) {
	// Create HTTPS test server
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("secure response"))
	}))
	defer server.Close()

	client := httpclient.New().
		WithTLSConfig(&tls.Config{
			InsecureSkipVerify: true, // For testing only
		})

	data, err := client.GET(server.URL)
	if err != nil {
		t.Fatalf("TLS request failed: %v", err)
	}

	if string(data) != "secure response" {
		t.Errorf("Unexpected TLS response: %s", data)
	}
}

func TestCookieJar(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/set-cookie" {
			http.SetCookie(w, &http.Cookie{
				Name:  "test-cookie",
				Value: "test-value",
			})
			w.Write([]byte("cookie set"))
		} else if r.URL.Path == "/check-cookie" {
			cookie, err := r.Cookie("test-cookie")
			if err != nil || cookie.Value != "test-value" {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("cookie not found"))
			} else {
				w.Write([]byte("cookie found"))
			}
		}
	}))
	defer server.Close()

	jar, _ := cookiejar.New(nil)
	client := httpclient.New().WithCookieJar(jar)

	// Set cookie
	_, err := client.GET(server.URL + "/set-cookie")
	if err != nil {
		t.Fatalf("Set cookie request failed: %v", err)
	}

	// Check cookie
	data, err := client.GET(server.URL + "/check-cookie")
	if err != nil {
		t.Fatalf("Check cookie request failed: %v", err)
	}

	if string(data) != "cookie found" {
		t.Errorf("Cookie not properly handled: %s", data)
	}
}

func TestRedirectPolicy(t *testing.T) {
	redirectCount := 0
	
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/redirect" {
			redirectCount++
			if redirectCount <= 2 {
				http.Redirect(w, r, "/redirect", http.StatusFound)
			} else {
				w.Write([]byte("final destination"))
			}
		}
	}))
	defer server.Close()

	client := httpclient.New().
		WithRedirectPolicy(func(req *http.Request, via []*http.Request) error {
			if len(via) >= 2 {
				return http.ErrUseLastResponse
			}
			return nil
		})

	_, err := client.GET(server.URL + "/redirect")
	if err != nil {
		t.Fatalf("Redirect policy request failed: %v", err)
	}

	if redirectCount > 2 {
		t.Errorf("Redirect policy not enforced, redirects: %d", redirectCount)
	}
}

func TestConnectionPool(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pooled connection"))
	}))
	defer server.Close()

	client := httpclient.New().
		WithConnectionPool(10, 5).
		WithKeepAlive(30 * time.Second)

	// Make multiple requests to test connection reuse
	for i := 0; i < 5; i++ {
		data, err := client.GET(server.URL)
		if err != nil {
			t.Fatalf("Connection pool request %d failed: %v", i, err)
		}
		if string(data) != "pooled connection" {
			t.Errorf("Unexpected response: %s", data)
		}
	}
}

func TestContextWithAdvancedFeatures(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond) // Simulate processing time
		w.Write([]byte("context response"))
	}))
	defer server.Close()

	client := httpclient.New().
		WithCompression(true).
		WithRequestInterceptor(func(req *http.Request) error {
			req.Header.Set("X-Context-Test", "true")
			return nil
		})

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	data, err := client.GetContext(ctx, server.URL)
	if err != nil {
		t.Fatalf("Context request with advanced features failed: %v", err)
	}

	if string(data) != "context response" {
		t.Errorf("Unexpected context response: %s", data)
	}
}