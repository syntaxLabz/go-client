package streaming

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// WebSocketConn represents a WebSocket connection
type WebSocketConn struct {
	conn   *websocket.Conn
	mu     sync.Mutex
	closed bool
}

// WebSocketDialer handles WebSocket connections
type WebSocketDialer struct {
	dialer   *websocket.Dialer
	headers  http.Header
	timeout  time.Duration
}

func NewWebSocketDialer() *WebSocketDialer {
	return &WebSocketDialer{
		dialer: &websocket.Dialer{
			HandshakeTimeout: 45 * time.Second,
			ReadBufferSize:   1024,
			WriteBufferSize:  1024,
		},
		headers: make(http.Header),
		timeout: 30 * time.Second,
	}
}

func (wd *WebSocketDialer) WithHeader(key, value string) *WebSocketDialer {
	wd.headers.Set(key, value)
	return wd
}

func (wd *WebSocketDialer) WithTimeout(timeout time.Duration) *WebSocketDialer {
	wd.timeout = timeout
	wd.dialer.HandshakeTimeout = timeout
	return wd
}

func (wd *WebSocketDialer) Dial(urlStr string) (*WebSocketConn, error) {
	return wd.DialContext(context.Background(), urlStr)
}

func (wd *WebSocketDialer) DialContext(ctx context.Context, urlStr string) (*WebSocketConn, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("invalid WebSocket URL: %w", err)
	}

	// Convert HTTP(S) URLs to WebSocket URLs
	switch u.Scheme {
	case "http":
		u.Scheme = "ws"
	case "https":
		u.Scheme = "wss"
	case "ws", "wss":
		// Already correct
	default:
		return nil, fmt.Errorf("unsupported scheme: %s", u.Scheme)
	}

	conn, _, err := wd.dialer.DialContext(ctx, u.String(), wd.headers)
	if err != nil {
		return nil, fmt.Errorf("WebSocket dial failed: %w", err)
	}

	return &WebSocketConn{
		conn: conn,
	}, nil
}

func (wc *WebSocketConn) Send(data interface{}) error {
	wc.mu.Lock()
	defer wc.mu.Unlock()

	if wc.closed {
		return fmt.Errorf("connection is closed")
	}

	var messageType int
	var payload []byte
	var err error

	switch v := data.(type) {
	case string:
		messageType = websocket.TextMessage
		payload = []byte(v)
	case []byte:
		messageType = websocket.BinaryMessage
		payload = v
	default:
		messageType = websocket.TextMessage
		payload, err = json.Marshal(data)
		if err != nil {
			return fmt.Errorf("failed to marshal data: %w", err)
		}
	}

	return wc.conn.WriteMessage(messageType, payload)
}

func (wc *WebSocketConn) Receive() ([]byte, error) {
	if wc.closed {
		return nil, fmt.Errorf("connection is closed")
	}

	_, data, err := wc.conn.ReadMessage()
	if err != nil {
		return nil, fmt.Errorf("failed to read message: %w", err)
	}

	return data, nil
}

func (wc *WebSocketConn) ReceiveJSON(v interface{}) error {
	data, err := wc.Receive()
	if err != nil {
		return err
	}

	return json.Unmarshal(data, v)
}

func (wc *WebSocketConn) Close() error {
	wc.mu.Lock()
	defer wc.mu.Unlock()

	if wc.closed {
		return nil
	}

	wc.closed = true
	return wc.conn.Close()
}

func (wc *WebSocketConn) SetReadDeadline(t time.Time) error {
	return wc.conn.SetReadDeadline(t)
}

func (wc *WebSocketConn) SetWriteDeadline(t time.Time) error {
	return wc.conn.SetWriteDeadline(t)
}

// StreamingClient handles streaming responses
type StreamingClient struct {
	client *http.Client
}

func NewStreamingClient() *StreamingClient {
	return &StreamingClient{
		client: &http.Client{
			Timeout: 0, // No timeout for streaming
		},
	}
}

func (sc *StreamingClient) Stream(method, url string, body interface{}) (<-chan []byte, error) {
	return sc.StreamContext(context.Background(), method, url, body)
}

func (sc *StreamingClient) StreamContext(ctx context.Context, method, url string, body interface{}) (<-chan []byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set streaming headers
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")

	resp, err := sc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.StatusCode >= 400 {
		resp.Body.Close()
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	ch := make(chan []byte, 100)

	go func() {
		defer resp.Body.Close()
		defer close(ch)

		buffer := make([]byte, 4096)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				n, err := resp.Body.Read(buffer)
				if err != nil {
					return
				}
				if n > 0 {
					data := make([]byte, n)
					copy(data, buffer[:n])
					select {
					case ch <- data:
					case <-ctx.Done():
						return
					}
				}
			}
		}
	}()

	return ch, nil
}

// ServerSentEvents handles SSE connections
type ServerSentEvents struct {
	client *http.Client
}

func NewServerSentEvents() *ServerSentEvents {
	return &ServerSentEvents{
		client: &http.Client{
			Timeout: 0, // No timeout for SSE
		},
	}
}

func (sse *ServerSentEvents) Connect(url string) (<-chan SSEEvent, error) {
	return sse.ConnectContext(context.Background(), url)
}

func (sse *ServerSentEvents) ConnectContext(ctx context.Context, url string) (<-chan SSEEvent, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")

	resp, err := sse.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.StatusCode >= 400 {
		resp.Body.Close()
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	ch := make(chan SSEEvent, 100)

	go func() {
		defer resp.Body.Close()
		defer close(ch)

		parser := NewSSEParser()
		buffer := make([]byte, 4096)

		for {
			select {
			case <-ctx.Done():
				return
			default:
				n, err := resp.Body.Read(buffer)
				if err != nil {
					return
				}
				if n > 0 {
					events := parser.Parse(buffer[:n])
					for _, event := range events {
						select {
						case ch <- event:
						case <-ctx.Done():
							return
						}
					}
				}
			}
		}
	}()

	return ch, nil
}

type SSEEvent struct {
	Type string
	Data string
	ID   string
}

type SSEParser struct {
	buffer []byte
}

func NewSSEParser() *SSEParser {
	return &SSEParser{
		buffer: make([]byte, 0),
	}
}

func (p *SSEParser) Parse(data []byte) []SSEEvent {
	p.buffer = append(p.buffer, data...)
	
	var events []SSEEvent
	lines := string(p.buffer)
	
	// Simple SSE parsing (production would need more robust parsing)
	if len(lines) > 0 {
		event := SSEEvent{
			Type: "message",
			Data: lines,
		}
		events = append(events, event)
		p.buffer = p.buffer[:0] // Clear buffer
	}
	
	return events
}