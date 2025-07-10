package main

import (
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
	fmt.Println("=== Developer Experience Features ===\n")

	// Example 1: Smart Constructors for Different Use Cases
	fmt.Println("1. Smart Constructors:")
	
	// Microservices client
	microserviceClient := httpclient.NewForMicroservices()
	fmt.Println("✓ Microservices client: Load balancing, health checks, circuit breaker, AI retry")
	
	// API client
	apiClient := httpclient.NewForAPI()
	fmt.Println("✓ API client: Rate limiting, caching, compression, smart caching")
	
	// Enterprise client
	enterpriseClient := httpclient.NewForEnterprise()
	fmt.Println("✓ Enterprise client: Full security, compliance, monitoring, AI features")
	
	// Development client
	devClient := httpclient.NewForDevelopment()
	fmt.Println("✓ Development client: Debugging, mocking, recording, chaos engineering")
	fmt.Println()

	// Example 2: One-liner Advanced Operations
	fmt.Println("2. One-liner Advanced Operations:")
	
	// Batch requests in one line
	fmt.Println("Batch requests:")
	responses, err := httpclient.Batch().
		Add("GET", "https://jsonplaceholder.typicode.com/users/1", nil).
		Add("GET", "https://jsonplaceholder.typicode.com/users/2", nil).
		Execute()
	
	if err != nil {
		log.Printf("Batch error: %v", err)
	} else {
		fmt.Printf("  Executed %d requests in batch\n", len(responses))
	}

	// GraphQL in one line
	fmt.Println("GraphQL query:")
	var result map[string]interface{}
	query := `{ __schema { queryType { name } } }`
	fmt.Printf("  Query: %s\n", query)
	fmt.Println("  (Would execute with valid GraphQL endpoint)")
	fmt.Println()

	// Example 3: Auto-configuration Based on Environment
	fmt.Println("3. Environment-aware Configuration:")
	
	// The client automatically detects environment and configures accordingly
	smartClient := httpclient.New().
		WithPerformanceOptimization(true).
		WithAIRetry(true).
		WithAdaptiveTimeout(true)

	fmt.Println("Smart client automatically configured with:")
	fmt.Println("  ✓ Performance optimization enabled")
	fmt.Println("  ✓ AI-powered retry logic")
	fmt.Println("  ✓ Adaptive timeout adjustment")
	fmt.Println("  ✓ Smart caching decisions")
	fmt.Println()

	// Example 4: Validation and Schema Support
	fmt.Println("4. Request/Response Validation:")
	
	validatingClient := httpclient.New().
		WithValidation(nil). // Would accept JSON schema
		WithDebug(true)

	var user User
	err = validatingClient.JSON("GET", "https://jsonplaceholder.typicode.com/users/1", nil, &user)
	if err != nil {
		log.Printf("Validation error: %v", err)
	} else {
		fmt.Printf("  Validated user: %+v\n", user)
	}
	fmt.Println()

	// Example 5: Chaos Engineering for Testing
	fmt.Println("5. Chaos Engineering:")
	
	chaosClient := httpclient.New().
		WithChaosEngineering(httpclient.ChaosConfig{
			FailureRate:    0.1,  // 10% failure rate
			LatencyRange:   [2]time.Duration{100*time.Millisecond, 500*time.Millisecond},
			ErrorTypes:     []string{"timeout", "connection_error", "server_error"},
			EnabledMethods: []string{"GET", "POST"},
		}).
		WithDebug(true)

	fmt.Println("Chaos engineering client configured:")
	fmt.Println("  ✓ 10% random failure injection")
	fmt.Println("  ✓ Random latency 100-500ms")
	fmt.Println("  ✓ Various error types")
	fmt.Println("  ✓ Enabled for GET/POST methods")
	fmt.Println()

	// Example 6: Recording and Replay for Testing
	fmt.Println("6. Recording and Replay:")
	
	recordingClient := httpclient.New().
		WithRecording(true).
		WithReplay(false). // Set to true to replay recorded responses
		WithDebug(true)

	fmt.Println("Recording client configured:")
	fmt.Println("  ✓ Records all requests and responses")
	fmt.Println("  ✓ Can replay recorded sessions")
	fmt.Println("  ✓ Perfect for testing and development")
	fmt.Println()

	// Example 7: Mocking Support
	fmt.Println("7. Built-in Mocking:")
	
	mockingClient := httpclient.New().
		WithMocking(true).
		WithDebug(true)

	fmt.Println("Mocking client configured:")
	fmt.Println("  ✓ Automatic mock responses for development")
	fmt.Println("  ✓ Configurable mock data")
	fmt.Println("  ✓ No external dependencies needed")
	fmt.Println()

	// Example 8: Complete Developer Experience
	fmt.Println("8. Complete Developer Experience:")
	
	completeClient := httpclient.New().
		WithAIRetry(true).
		WithSmartCaching(true).
		WithAdaptiveTimeout(true).
		WithPerformanceOptimization(true).
		WithRealTimeMetrics(true).
		WithValidation(nil).
		WithDebug(true).
		WithAutoRetry(httpclient.AutoRetryConfig{
			MaxAttempts:     5,
			BackoffStrategy: "exponential",
			RetryConditions: []string{"timeout", "5xx", "connection_error"},
			JitterEnabled:   true,
		})

	fmt.Println("Complete developer experience client:")
	fmt.Println("  ✓ AI-powered features")
	fmt.Println("  ✓ Smart caching and preloading")
	fmt.Println("  ✓ Adaptive performance optimization")
	fmt.Println("  ✓ Real-time metrics and monitoring")
	fmt.Println("  ✓ Request/response validation")
	fmt.Println("  ✓ Comprehensive debugging")
	fmt.Println("  ✓ Intelligent auto-retry")
	fmt.Println("  ✓ Jitter for retry timing")

	// Test the complete client
	var testUser User
	err = completeClient.JSON("GET", "https://jsonplaceholder.typicode.com/users/1", nil, &testUser)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("  Result: %s (%s)\n", testUser.Name, testUser.Email)
	}
	fmt.Println()

	// Example 9: Ultra-Simple Usage
	fmt.Println("9. Ultra-Simple Usage Examples:")
	fmt.Println()
	
	fmt.Println("// Basic requests")
	fmt.Println("data, err := httpclient.GET(\"https://api.example.com/users\")")
	fmt.Println("err := httpclient.JSON(\"POST\", \"https://api.example.com/users\", user, &result)")
	fmt.Println()
	
	fmt.Println("// Advanced features in one line")
	fmt.Println("client := httpclient.NewForMicroservices()")
	fmt.Println("responses, err := httpclient.Batch().Add(...).Add(...).Execute()")
	fmt.Println("ws, err := httpclient.WebSocket(\"wss://api.example.com\")")
	fmt.Println("err := httpclient.GraphQL(query, variables, &result)")
	fmt.Println()
	
	fmt.Println("// AI-powered features")
	fmt.Println("client := httpclient.New().WithAIRetry(true).WithSmartCaching(true)")
	fmt.Println("stream, err := httpclient.Stream(\"GET\", \"https://api.example.com/stream\", nil)")
	fmt.Println()
	
	fmt.Println("The client is designed to be:")
	fmt.Println("  🚀 Ultra-simple for basic use cases")
	fmt.Println("  🧠 AI-powered for intelligent behavior")
	fmt.Println("  🏗️  Enterprise-ready with advanced features")
	fmt.Println("  🔧 Developer-friendly with great tooling")
	fmt.Println("  📊 Observable with comprehensive metrics")
	fmt.Println("  🛡️  Secure with built-in best practices")
}