package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/yourorg/httpclient"
)

func main() {
	fmt.Println("=== Streaming and Real-time Features ===\n")

	// Example 1: HTTP Streaming
	fmt.Println("1. HTTP Streaming:")
	streamingClient := httpclient.New().
		WithTimeout(0). // No timeout for streaming
		WithDebug(true)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Note: This would work with a real streaming endpoint
	fmt.Println("Streaming client configured (would work with real streaming endpoints)")
	fmt.Println("Example usage:")
	fmt.Println(`
	stream, err := streamingClient.Stream("GET", "https://api.example.com/stream", nil)
	if err != nil {
		log.Fatal(err)
	}
	
	for data := range stream {
		fmt.Printf("Received: %s\n", data)
	}
	`)
	fmt.Println()

	// Example 2: WebSocket Connection
	fmt.Println("2. WebSocket Support:")
	wsClient := httpclient.New().
		WithTimeout(30 * time.Second)

	fmt.Println("WebSocket client configured")
	fmt.Println("Example usage:")
	fmt.Println(`
	ws, err := wsClient.WebSocket("wss://echo.websocket.org")
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()
	
	// Send message
	err = ws.Send("Hello WebSocket!")
	if err != nil {
		log.Fatal(err)
	}
	
	// Receive message
	data, err := ws.Receive()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Received: %s\n", data)
	`)
	fmt.Println()

	// Example 3: Batch Requests
	fmt.Println("3. Batch Request Processing:")
	batchClient := httpclient.New().
		WithTimeout(30 * time.Second).
		WithDebug(true)

	// Create a batch of requests
	batch := batchClient.Batch().
		Add("GET", "https://jsonplaceholder.typicode.com/users/1", nil).
		Add("GET", "https://jsonplaceholder.typicode.com/users/2", nil).
		Add("GET", "https://jsonplaceholder.typicode.com/users/3", nil).
		Add("GET", "https://jsonplaceholder.typicode.com/posts/1", nil)

	fmt.Println("Executing batch of 4 requests concurrently...")
	start := time.Now()
	responses, err := batch.Execute()
	duration := time.Since(start)

	if err != nil {
		log.Printf("Batch error: %v", err)
	} else {
		fmt.Printf("Batch completed in %v\n", duration)
		for _, resp := range responses {
			if resp.Error != nil {
				fmt.Printf("  Request %d failed: %v\n", resp.Index, resp.Error)
			} else {
				fmt.Printf("  Request %d succeeded in %v (size: %d bytes)\n", 
					resp.Index, resp.Duration, len(resp.Data))
			}
		}
	}
	fmt.Println()

	// Example 4: Pipeline Requests
	fmt.Println("4. Pipeline Request Processing:")
	pipelineClient := httpclient.New().
		WithTimeout(30 * time.Second)

	// Create a pipeline of requests
	pipeline := pipelineClient.Pipeline().
		Add("GET", "https://jsonplaceholder.typicode.com/users/1", nil).
		Add("GET", "https://jsonplaceholder.typicode.com/users/2", nil).
		Add("GET", "https://jsonplaceholder.typicode.com/users/3", nil)

	fmt.Println("Executing pipeline of 3 requests sequentially...")
	responseChan, err := pipeline.Execute()
	if err != nil {
		log.Printf("Pipeline error: %v", err)
	} else {
		for resp := range responseChan {
			if resp.Error != nil {
				fmt.Printf("  Pipeline request %d failed: %v\n", resp.Index, resp.Error)
			} else {
				fmt.Printf("  Pipeline request %d completed in %v (size: %d bytes)\n", 
					resp.Index, resp.Duration, len(resp.Data))
			}
		}
	}
	fmt.Println()

	// Example 5: GraphQL Support
	fmt.Println("5. GraphQL Support:")
	graphqlClient := httpclient.New().
		WithBaseURL("https://api.github.com/graphql").
		WithHeader("Authorization", "Bearer YOUR_TOKEN_HERE").
		WithTimeout(30 * time.Second)

	query := `
		query {
			viewer {
				login
				name
				email
			}
		}
	`

	fmt.Println("GraphQL client configured")
	fmt.Println("Example query:")
	fmt.Println(query)
	fmt.Println("Note: Requires valid GitHub token for actual execution")
	fmt.Println()

	// Example 6: Real-time Features Combined
	fmt.Println("6. Complete Real-time Configuration:")
	realtimeClient := httpclient.New().
		WithTimeout(0). // No timeout for real-time
		WithRealTimeMetrics(true).
		WithPerformanceOptimization(true).
		WithDebug(true)

	fmt.Println("Real-time client configured with:")
	fmt.Println("  ✓ HTTP streaming support")
	fmt.Println("  ✓ WebSocket connections")
	fmt.Println("  ✓ Server-Sent Events (SSE)")
	fmt.Println("  ✓ Batch request processing")
	fmt.Println("  ✓ Pipeline request processing")
	fmt.Println("  ✓ GraphQL support")
	fmt.Println("  ✓ Real-time metrics")
	fmt.Println("  ✓ Performance optimization")

	// Demonstrate simple usage
	fmt.Println("\nSimple usage examples:")
	fmt.Println("// Batch requests")
	fmt.Println("responses, err := httpclient.Batch().Add(...).Add(...).Execute()")
	fmt.Println()
	fmt.Println("// Pipeline requests") 
	fmt.Println("stream, err := httpclient.Pipeline().Add(...).Add(...).Execute()")
	fmt.Println()
	fmt.Println("// WebSocket")
	fmt.Println("ws, err := httpclient.WebSocket(\"wss://example.com\")")
	fmt.Println()
	fmt.Println("// GraphQL")
	fmt.Println("err := httpclient.GraphQL(query, variables, &result)")
	fmt.Println()
	fmt.Println("// HTTP Streaming")
	fmt.Println("stream, err := httpclient.Stream(\"GET\", \"https://api.example.com/stream\", nil)")
}