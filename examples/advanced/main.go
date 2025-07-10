package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/yourorg/httpclient"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	fmt.Println("=== Advanced HTTP Client Examples ===\n")

	// Example 1: Custom client with configuration
	fmt.Println("1. Custom client with configuration:")
	client := httpclient.New().
		WithTimeout(10 * time.Second).
		WithRetries(5).
		WithUserAgent("MyApp/1.0").
		WithHeader("X-API-Version", "v1").
		WithDebug(true)

	var user User
	err := client.JSON("GET", "https://jsonplaceholder.typicode.com/users/1", nil, &user)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("User: %+v\n\n", user)
	}

	// Example 2: Authentication
	fmt.Println("2. Client with authentication:")
	authClient := httpclient.New().
		WithAuth("your-bearer-token").
		WithAPIKey("X-API-Key", "your-api-key")

	// This would work with a real authenticated API
	fmt.Println("Auth client configured (would work with real authenticated API)\n")

	// Example 3: Base URL for API clients
	fmt.Println("3. API client with base URL:")
	apiClient := httpclient.New().
		WithBaseURL("https://jsonplaceholder.typicode.com").
		WithTimeout(5 * time.Second)

	var users []User
	err = apiClient.JSON("GET", "/users", nil, &users)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("Found %d users\n\n", len(users))
	}

	// Example 4: Context-aware requests
	fmt.Println("4. Context-aware requests:")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	data, err := httpclient.GetContext(ctx, "https://jsonplaceholder.typicode.com/users/1")
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("Response length: %d bytes\n\n", len(data))
	}

	// Example 5: Rate limiting and circuit breaker
	fmt.Println("5. Advanced features:")
	advancedClient := httpclient.New().
		WithRateLimiter(10).                                    // 10 requests per second
		WithCircuitBreaker(3, 30*time.Second).                 // Break after 3 failures for 30s
		WithCache(5 * time.Minute).                            // Cache responses for 5 minutes
		WithMetrics(true).                                      // Enable Prometheus metrics
		WithTracing(true)                                       // Enable OpenTelemetry tracing

	fmt.Println("Advanced client configured with rate limiting, circuit breaker, caching, metrics, and tracing\n")

	// Example 6: Fluent interface chaining
	fmt.Println("6. Fluent interface chaining:")
	response, err := httpclient.New().
		WithBaseURL("https://jsonplaceholder.typicode.com").
		WithTimeout(10 * time.Second).
		WithRetries(3).
		WithUserAgent("FluentClient/1.0").
		WithHeader("Accept", "application/json").
		GET("/users/1")

	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("Fluent response length: %d bytes\n", len(response))
	}
}