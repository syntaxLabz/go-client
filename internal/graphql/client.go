package graphql

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// GraphQLClient handles GraphQL requests
type GraphQLClient struct {
	endpoint string
	client   *http.Client
	headers  map[string]string
}

type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

type GraphQLResponse struct {
	Data   json.RawMessage `json:"data"`
	Errors []GraphQLError  `json:"errors,omitempty"`
}

type GraphQLError struct {
	Message   string                 `json:"message"`
	Locations []GraphQLLocation      `json:"locations,omitempty"`
	Path      []interface{}          `json:"path,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
}

type GraphQLLocation struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

func NewGraphQLClient(endpoint string, client *http.Client) *GraphQLClient {
	return &GraphQLClient{
		endpoint: endpoint,
		client:   client,
		headers:  make(map[string]string),
	}
}

func (gc *GraphQLClient) WithHeader(key, value string) *GraphQLClient {
	gc.headers[key] = value
	return gc
}

func (gc *GraphQLClient) Query(query string, variables map[string]interface{}, result interface{}) error {
	return gc.QueryContext(context.Background(), query, variables, result)
}

func (gc *GraphQLClient) QueryContext(ctx context.Context, query string, variables map[string]interface{}, result interface{}) error {
	reqBody := GraphQLRequest{
		Query:     query,
		Variables: variables,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal GraphQL request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", gc.endpoint, bytes.NewReader(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Add custom headers
	for key, value := range gc.headers {
		req.Header.Set(key, value)
	}

	resp, err := gc.client.Do(req)
	if err != nil {
		return fmt.Errorf("GraphQL request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("GraphQL HTTP error: %d %s", resp.StatusCode, resp.Status)
	}

	var gqlResp GraphQLResponse
	if err := json.NewDecoder(resp.Body).Decode(&gqlResp); err != nil {
		return fmt.Errorf("failed to decode GraphQL response: %w", err)
	}

	if len(gqlResp.Errors) > 0 {
		return &GraphQLErrors{Errors: gqlResp.Errors}
	}

	if result != nil && len(gqlResp.Data) > 0 {
		if err := json.Unmarshal(gqlResp.Data, result); err != nil {
			return fmt.Errorf("failed to unmarshal GraphQL data: %w", err)
		}
	}

	return nil
}

// GraphQLErrors represents multiple GraphQL errors
type GraphQLErrors struct {
	Errors []GraphQLError
}

func (e *GraphQLErrors) Error() string {
	if len(e.Errors) == 1 {
		return fmt.Sprintf("GraphQL error: %s", e.Errors[0].Message)
	}
	return fmt.Sprintf("GraphQL errors: %d errors occurred", len(e.Errors))
}

// Subscription support for GraphQL subscriptions over WebSocket
type GraphQLSubscription struct {
	client   *GraphQLClient
	wsClient interface{} // WebSocket client would be injected
}

func (gc *GraphQLClient) Subscribe(query string, variables map[string]interface{}) (*GraphQLSubscription, error) {
	// This would implement GraphQL subscriptions over WebSocket
	// For now, return a placeholder
	return &GraphQLSubscription{
		client: gc,
	}, nil
}

// Helper functions for common GraphQL operations
func (gc *GraphQLClient) Introspect() (map[string]interface{}, error) {
	introspectionQuery := `
		query IntrospectionQuery {
			__schema {
				types {
					name
					kind
					description
					fields {
						name
						type {
							name
							kind
						}
					}
				}
			}
		}
	`

	var result map[string]interface{}
	err := gc.Query(introspectionQuery, nil, &result)
	return result, err
}

func (gc *GraphQLClient) GetSchema() (string, error) {
	schemaQuery := `
		query GetSchema {
			__schema {
				queryType { name }
				mutationType { name }
				subscriptionType { name }
			}
		}
	`

	var result map[string]interface{}
	err := gc.Query(schemaQuery, nil, &result)
	if err != nil {
		return "", err
	}

	schema, _ := json.MarshalIndent(result, "", "  ")
	return string(schema), nil
}