package test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourorg/httpclient"
)

type TestUser struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func TestBasicHTTPMethods(t *testing.T) {
	tests := []struct {
		name   string
		method string
		fn     func(string) ([]byte, error)
	}{
		{"GET", "GET", httpclient.GET},
		{"DELETE", "DELETE", httpclient.DELETE},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != tt.method {
					t.Errorf("Expected %s, got %s", tt.method, r.Method)
				}
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"id": 1, "name": "John"}`))
			}))
			defer server.Close()

			data, err := tt.fn(server.URL)
			if err != nil {
				t.Fatalf("%s failed: %v", tt.method, err)
			}

			var user TestUser
			if err := json.Unmarshal(data, &user); err != nil {
				t.Fatalf("JSON unmarshal failed: %v", err)
			}

			if user.ID != 1 || user.Name != "John" {
				t.Errorf("Expected user {1, John}, got %+v", user)
			}
		})
	}
}

func TestPOSTWithBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type application/json")
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id": 2, "name": "Jane"}`))
	}))
	defer server.Close()

	user := TestUser{Name: "Jane"}
	data, err := httpclient.POST(server.URL, user)
	if err != nil {
		t.Fatalf("POST failed: %v", err)
	}

	var result TestUser
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("JSON unmarshal failed: %v", err)
	}

	if result.ID != 2 || result.Name != "Jane" {
		t.Errorf("Expected user {2, Jane}, got %+v", result)
	}
}

func TestJSONMethod(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": 3, "name": "Bob"}`))
	}))
	defer server.Close()

	var user TestUser
	err := httpclient.JSON("GET", server.URL, nil, &user)
	if err != nil {
		t.Fatalf("JSON failed: %v", err)
	}

	if user.ID != 3 || user.Name != "Bob" {
		t.Errorf("Expected user {3, Bob}, got %+v", user)
	}
}

func TestClientConfiguration(t *testing.T) {
	client := httpclient.New().
		WithTimeout(5 * time.Second).
		WithRetries(2).
		WithAuth("test-token").
		WithHeader("X-Test", "value").
		WithUserAgent("TestClient/1.0")

	// Test that client is properly configured (this is a basic test)
	// In a real scenario, you'd test the actual behavior
	if client == nil {
		t.Error("Client should not be nil")
	}
}

func TestContextAwareMethods(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "ok"}`))
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	data, err := httpclient.GetContext(ctx, server.URL)
	if err != nil {
		t.Fatalf("GetContext failed: %v", err)
	}

	if string(data) != `{"message": "ok"}` {
		t.Errorf("Unexpected response: %s", data)
	}
}

func TestErrorHandling(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not found"))
	}))
	defer server.Close()

	_, err := httpclient.GET(server.URL)
	if err == nil {
		t.Fatal("Expected error for 404 response")
	}

	if err.Error() != "HTTP 404: Not found" {
		t.Errorf("Unexpected error message: %v", err)
	}
}

func TestFluentInterface(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check headers set by fluent interface
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("Expected auth header")
		}
		if r.Header.Get("X-Custom") != "custom-value" {
			t.Errorf("Expected custom header")
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true}`))
	}))
	defer server.Close()

	data, err := httpclient.New().
		WithAuth("test-token").
		WithHeader("X-Custom", "custom-value").
		WithTimeout(10 * time.Second).
		GET(server.URL)

	if err != nil {
		t.Fatalf("Fluent interface request failed: %v", err)
	}

	if string(data) != `{"success": true}` {
		t.Errorf("Unexpected response: %s", data)
	}
}

func TestBaseURL(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/users" {
			t.Errorf("Expected path /api/users, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"users": []}`))
	}))
	defer server.Close()

	client := httpclient.New().WithBaseURL(server.URL)
	data, err := client.GET("/api/users")
	if err != nil {
		t.Fatalf("BaseURL request failed: %v", err)
	}

	if string(data) != `{"users": []}` {
		t.Errorf("Unexpected response: %s", data)
	}
}