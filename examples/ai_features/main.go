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
	fmt.Println("=== AI-Powered HTTP Client Features ===\n")

	// Example 1: AI-Enhanced Retry Strategy
	fmt.Println("1. AI-Enhanced Smart Retry:")
	smartClient := httpclient.New().
		WithAIRetry(true).
		WithAdaptiveTimeout(true).
		WithTimeout(10 * time.Second)

	var user User
	err := smartClient.JSON("GET", "https://jsonplaceholder.typicode.com/users/1", nil, &user)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("User with AI retry: %+v\n\n", user)
	}

	// Example 2: Smart Caching with AI
	fmt.Println("2. AI-Powered Smart Caching:")
	cachingClient := httpclient.New().
		WithSmartCaching(true).
		WithCache(5 * time.Minute).
		WithPredictivePreloading(true)

	// Make multiple requests to see smart caching in action
	for i := 1; i <= 3; i++ {
		var user User
		start := time.Now()
		err := cachingClient.JSON("GET", fmt.Sprintf("https://jsonplaceholder.typicode.com/users/%d", i), nil, &user)
		duration := time.Since(start)
		
		if err != nil {
			log.Printf("Error: %v", err)
		} else {
			fmt.Printf("User %d fetched in %v: %s\n", i, duration, user.Name)
		}
	}
	fmt.Println()

	// Example 3: Adaptive Timeout Based on Performance
	fmt.Println("3. Adaptive Timeout Optimization:")
	adaptiveClient := httpclient.New().
		WithAdaptiveTimeout(true).
		WithPerformanceOptimization(true)

	// Simulate multiple requests to different endpoints
	endpoints := []string{
		"https://jsonplaceholder.typicode.com/users/1",
		"https://jsonplaceholder.typicode.com/posts/1",
		"https://jsonplaceholder.typicode.com/comments/1",
	}

	for _, endpoint := range endpoints {
		start := time.Now()
		data, err := adaptiveClient.GET(endpoint)
		duration := time.Since(start)
		
		if err != nil {
			log.Printf("Error for %s: %v", endpoint, err)
		} else {
			fmt.Printf("Endpoint %s responded in %v (size: %d bytes)\n", endpoint, duration, len(data))
		}
	}
	fmt.Println()

	// Example 4: Predictive Preloading
	fmt.Println("4. Predictive Preloading:")
	preloadingClient := httpclient.New().
		WithPredictivePreloading(true).
		WithSmartCaching(true)

	// Simulate a pattern of requests
	fmt.Println("Simulating user browsing pattern...")
	
	// First request - user profile
	var userProfile User
	err = preloadingClient.JSON("GET", "https://jsonplaceholder.typicode.com/users/1", nil, &userProfile)
	if err == nil {
		fmt.Printf("Loaded user profile: %s\n", userProfile.Name)
	}

	// AI should predict that posts might be requested next
	time.Sleep(100 * time.Millisecond) // Simulate user interaction

	// Second request - user's posts (should be faster due to preloading)
	start := time.Now()
	posts, err := preloadingClient.GET("https://jsonplaceholder.typicode.com/users/1/posts")
	duration := time.Since(start)
	
	if err == nil {
		fmt.Printf("Loaded user posts in %v (size: %d bytes)\n", duration, len(posts))
	}
	fmt.Println()

	// Example 5: Complete AI-Powered Client
	fmt.Println("5. Complete AI-Powered Configuration:")
	aiClient := httpclient.New().
		WithAIRetry(true).
		WithSmartCaching(true).
		WithPredictivePreloading(true).
		WithAdaptiveTimeout(true).
		WithPerformanceOptimization(true).
		WithRealTimeMetrics(true).
		WithDebug(true)

	fmt.Println("AI client configured with:")
	fmt.Println("  ✓ Smart retry with machine learning")
	fmt.Println("  ✓ Intelligent caching decisions")
	fmt.Println("  ✓ Predictive request preloading")
	fmt.Println("  ✓ Adaptive timeout optimization")
	fmt.Println("  ✓ Real-time performance optimization")
	fmt.Println("  ✓ Live metrics and monitoring")

	// Test the AI client
	var testUser User
	err = aiClient.JSON("GET", "https://jsonplaceholder.typicode.com/users/2", nil, &testUser)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("AI client result: %+v\n", testUser)
	}
}