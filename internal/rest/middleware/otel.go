package middleware

import (
	"time"

	"github.com/edwinjordan/ZOGTest-Golang.git/internal/metrics"
	echo "github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// OTelMetricsMiddleware creates middleware for recording OpenTelemetry metrics
func OTelMetricsMiddleware(otelMetrics *metrics.OTelMetrics) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()
			start := time.Now()

			// Track in-flight requests
			otelMetrics.HTTPRequestsInFlight.Add(ctx, 1)
			defer otelMetrics.HTTPRequestsInFlight.Add(ctx, -1)

			// Process request
			err := next(c)

			// Calculate duration in milliseconds
			duration := float64(time.Since(start).Milliseconds())

			// Record metrics
			req := c.Request()
			res := c.Response()
			otelMetrics.RecordHTTPRequest(
				ctx,
				req.Method,
				c.Path(), // Use route pattern, not actual path
				res.Status,
				duration,
			)

			return err
		}
	}
}

// EnhancedTracingMiddleware adds detailed tracing with spans
func EnhancedTracingMiddleware(serviceName string) echo.MiddlewareFunc {
	tracer := otel.Tracer(serviceName)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			ctx := req.Context()

			// Start a new span
			spanName := req.Method + " " + c.Path()
			ctx, span := tracer.Start(ctx, spanName,
				trace.WithSpanKind(trace.SpanKindServer),
				trace.WithAttributes(
					attribute.String("http.method", req.Method),
					attribute.String("http.url", req.URL.String()),
					attribute.String("http.scheme", req.URL.Scheme),
					attribute.String("http.target", req.URL.Path),
					attribute.String("http.route", c.Path()),
					attribute.String("http.user_agent", req.UserAgent()),
					attribute.String("http.client_ip", c.RealIP()),
					attribute.String("net.host.name", req.Host),
				),
			)
			defer span.End()

			// Update request context
			c.SetRequest(req.WithContext(ctx))

			// Process request
			err := next(c)

			// Record response details
			res := c.Response()
			span.SetAttributes(
				attribute.Int("http.status_code", res.Status),
				attribute.Int64("http.response_content_length", res.Size),
			)

			// Set span status based on HTTP status code
			if res.Status >= 500 {
				span.SetStatus(codes.Error, "Internal Server Error")
			} else if res.Status >= 400 {
				span.SetStatus(codes.Error, "Client Error")
			} else {
				span.SetStatus(codes.Ok, "Success")
			}

			// Record error if present
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			}

			return err
		}
	}
}
