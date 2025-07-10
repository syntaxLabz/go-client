package batch

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// BatchRequest represents a batch of HTTP requests
type BatchRequest struct {
	requests []BatchItem
	client   *http.Client
	mu       sync.Mutex
}

type BatchItem struct {
	Method string
	URL    string
	Body   interface{}
	Index  int
}

type BatchResponse struct {
	Index    int
	Data     []byte
	Error    error
	Duration time.Duration
}

func NewBatchRequest(client *http.Client) *BatchRequest {
	return &BatchRequest{
		requests: make([]BatchItem, 0),
		client:   client,
	}
}

func (br *BatchRequest) Add(method, url string, body interface{}) *BatchRequest {
	br.mu.Lock()
	defer br.mu.Unlock()
	
	br.requests = append(br.requests, BatchItem{
		Method: method,
		URL:    url,
		Body:   body,
		Index:  len(br.requests),
	})
	
	return br
}

func (br *BatchRequest) Execute() ([]BatchResponse, error) {
	return br.ExecuteContext(context.Background())
}

func (br *BatchRequest) ExecuteContext(ctx context.Context) ([]BatchResponse, error) {
	br.mu.Lock()
	requests := make([]BatchItem, len(br.requests))
	copy(requests, br.requests)
	br.mu.Unlock()

	responses := make([]BatchResponse, len(requests))
	var wg sync.WaitGroup
	
	// Execute all requests concurrently
	for i, req := range requests {
		wg.Add(1)
		go func(index int, item BatchItem) {
			defer wg.Done()
			
			start := time.Now()
			data, err := br.executeRequest(ctx, item)
			duration := time.Since(start)
			
			responses[index] = BatchResponse{
				Index:    item.Index,
				Data:     data,
				Error:    err,
				Duration: duration,
			}
		}(i, req)
	}
	
	wg.Wait()
	return responses, nil
}

func (br *BatchRequest) executeRequest(ctx context.Context, item BatchItem) ([]byte, error) {
	var reqBody []byte
	var err error
	
	if item.Body != nil {
		reqBody, err = json.Marshal(item.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal body: %w", err)
		}
	}
	
	req, err := http.NewRequestWithContext(ctx, item.Method, item.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	if reqBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	
	resp, err := br.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}
	
	// Read response body
	data := make([]byte, 0)
	buffer := make([]byte, 4096)
	for {
		n, err := resp.Body.Read(buffer)
		if n > 0 {
			data = append(data, buffer[:n]...)
		}
		if err != nil {
			break
		}
	}
	
	return data, nil
}

// PipelineRequest represents a pipeline of HTTP requests
type PipelineRequest struct {
	requests []BatchItem
	client   *http.Client
	mu       sync.Mutex
}

type PipelineResponse struct {
	Index    int
	Data     []byte
	Error    error
	Duration time.Duration
}

func NewPipelineRequest(client *http.Client) *PipelineRequest {
	return &PipelineRequest{
		requests: make([]BatchItem, 0),
		client:   client,
	}
}

func (pr *PipelineRequest) Add(method, url string, body interface{}) *PipelineRequest {
	pr.mu.Lock()
	defer pr.mu.Unlock()
	
	pr.requests = append(pr.requests, BatchItem{
		Method: method,
		URL:    url,
		Body:   body,
		Index:  len(pr.requests),
	})
	
	return pr
}

func (pr *PipelineRequest) Execute() (<-chan PipelineResponse, error) {
	return pr.ExecuteContext(context.Background())
}

func (pr *PipelineRequest) ExecuteContext(ctx context.Context) (<-chan PipelineResponse, error) {
	pr.mu.Lock()
	requests := make([]BatchItem, len(pr.requests))
	copy(requests, pr.requests)
	pr.mu.Unlock()

	ch := make(chan PipelineResponse, len(requests))
	
	go func() {
		defer close(ch)
		
		// Execute requests in sequence, streaming results
		for _, req := range requests {
			start := time.Now()
			data, err := pr.executeRequest(ctx, req)
			duration := time.Since(start)
			
			response := PipelineResponse{
				Index:    req.Index,
				Data:     data,
				Error:    err,
				Duration: duration,
			}
			
			select {
			case ch <- response:
			case <-ctx.Done():
				return
			}
		}
	}()
	
	return ch, nil
}

func (pr *PipelineRequest) executeRequest(ctx context.Context, item BatchItem) ([]byte, error) {
	var reqBody []byte
	var err error
	
	if item.Body != nil {
		reqBody, err = json.Marshal(item.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal body: %w", err)
		}
	}
	
	req, err := http.NewRequestWithContext(ctx, item.Method, item.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	if reqBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	
	resp, err := pr.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}
	
	// Read response body
	data := make([]byte, 0)
	buffer := make([]byte, 4096)
	for {
		n, err := resp.Body.Read(buffer)
		if n > 0 {
			data = append(data, buffer[:n]...)
		}
		if err != nil {
			break
		}
	}
	
	return data, nil
}