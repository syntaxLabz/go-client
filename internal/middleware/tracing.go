package middleware

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type tracingMiddleware struct {
	tracer trace.Tracer
	span   trace.Span
}

// NewTracing creates a new tracing middleware
func NewTracing() Middleware {
	return &tracingMiddleware{
		tracer: otel.Tracer("httpclient"),
	}
}

func (t *tracingMiddleware) Before(req *http.Request) error {
	ctx := req.Context()
	ctx, span := t.tracer.Start(ctx, "http_request",
		trace.WithAttributes(
			attribute.String("http.method", req.Method),
			attribute.String("http.url", req.URL.String()),
		),
	)
	
	t.span = span
	*req = *req.WithContext(ctx)
	return nil
}

func (t *tracingMiddleware) After(resp *http.Response) {
	if t.span != nil {
		t.span.SetAttributes(
			attribute.Int("http.status_code", resp.StatusCode),
		)
		t.span.End()
	}
}