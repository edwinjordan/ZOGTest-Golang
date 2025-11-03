package metrics

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// OTelMetrics holds OpenTelemetry metric instruments
type OTelMetrics struct {
	// HTTP metrics
	HTTPRequestsTotal    metric.Int64Counter
	HTTPRequestDuration  metric.Float64Histogram
	HTTPRequestsInFlight metric.Int64UpDownCounter

	// Database metrics
	DBQueriesTotal    metric.Int64Counter
	DBQueryDuration   metric.Float64Histogram
	DBConnectionsOpen metric.Int64UpDownCounter

	// Business metrics
	UserOperationsTotal  metric.Int64Counter
	TopicOperationsTotal metric.Int64Counter
	NewsOperationsTotal  metric.Int64Counter
}

// NewOTelMetrics creates a new set of OpenTelemetry metrics
func NewOTelMetrics(serviceName string) (*OTelMetrics, error) {
	meter := otel.Meter(serviceName)

	// HTTP metrics
	httpRequestsTotal, err := meter.Int64Counter(
		"http.server.requests.total",
		metric.WithDescription("Total number of HTTP requests"),
		metric.WithUnit("{request}"),
	)
	if err != nil {
		return nil, err
	}

	httpRequestDuration, err := meter.Float64Histogram(
		"http.server.request.duration",
		metric.WithDescription("HTTP request duration in milliseconds"),
		metric.WithUnit("ms"),
	)
	if err != nil {
		return nil, err
	}

	httpRequestsInFlight, err := meter.Int64UpDownCounter(
		"http.server.requests.in_flight",
		metric.WithDescription("Number of HTTP requests currently in flight"),
		metric.WithUnit("{request}"),
	)
	if err != nil {
		return nil, err
	}

	// Database metrics
	dbQueriesTotal, err := meter.Int64Counter(
		"db.queries.total",
		metric.WithDescription("Total number of database queries"),
		metric.WithUnit("{query}"),
	)
	if err != nil {
		return nil, err
	}

	dbQueryDuration, err := meter.Float64Histogram(
		"db.query.duration",
		metric.WithDescription("Database query duration in milliseconds"),
		metric.WithUnit("ms"),
	)
	if err != nil {
		return nil, err
	}

	dbConnectionsOpen, err := meter.Int64UpDownCounter(
		"db.connections.open",
		metric.WithDescription("Number of open database connections"),
		metric.WithUnit("{connection}"),
	)
	if err != nil {
		return nil, err
	}

	// Business metrics
	userOperationsTotal, err := meter.Int64Counter(
		"business.user.operations.total",
		metric.WithDescription("Total number of user operations"),
		metric.WithUnit("{operation}"),
	)
	if err != nil {
		return nil, err
	}

	topicOperationsTotal, err := meter.Int64Counter(
		"business.topic.operations.total",
		metric.WithDescription("Total number of topic operations"),
		metric.WithUnit("{operation}"),
	)
	if err != nil {
		return nil, err
	}

	newsOperationsTotal, err := meter.Int64Counter(
		"business.news.operations.total",
		metric.WithDescription("Total number of news operations"),
		metric.WithUnit("{operation}"),
	)
	if err != nil {
		return nil, err
	}

	slog.Info("OpenTelemetry metrics instruments created successfully")

	return &OTelMetrics{
		HTTPRequestsTotal:    httpRequestsTotal,
		HTTPRequestDuration:  httpRequestDuration,
		HTTPRequestsInFlight: httpRequestsInFlight,
		DBQueriesTotal:       dbQueriesTotal,
		DBQueryDuration:      dbQueryDuration,
		DBConnectionsOpen:    dbConnectionsOpen,
		UserOperationsTotal:  userOperationsTotal,
		TopicOperationsTotal: topicOperationsTotal,
		NewsOperationsTotal:  newsOperationsTotal,
	}, nil
}

// RecordHTTPRequest records an HTTP request metric
func (m *OTelMetrics) RecordHTTPRequest(ctx context.Context, method, path string, statusCode int, duration float64) {
	attrs := []attribute.KeyValue{
		attribute.String("http.method", method),
		attribute.String("http.route", path),
		attribute.Int("http.status_code", statusCode),
	}

	m.HTTPRequestsTotal.Add(ctx, 1, metric.WithAttributes(attrs...))
	m.HTTPRequestDuration.Record(ctx, duration, metric.WithAttributes(attrs...))
}

// RecordDBQuery records a database query metric
func (m *OTelMetrics) RecordDBQuery(ctx context.Context, operation, table string, duration float64, success bool) {
	attrs := []attribute.KeyValue{
		attribute.String("db.operation", operation),
		attribute.String("db.table", table),
		attribute.Bool("success", success),
	}

	m.DBQueriesTotal.Add(ctx, 1, metric.WithAttributes(attrs...))
	m.DBQueryDuration.Record(ctx, duration, metric.WithAttributes(attrs...))
}

// RecordUserOperation records a user operation metric
func (m *OTelMetrics) RecordUserOperation(ctx context.Context, operation string, success bool) {
	attrs := []attribute.KeyValue{
		attribute.String("operation", operation),
		attribute.Bool("success", success),
	}
	m.UserOperationsTotal.Add(ctx, 1, metric.WithAttributes(attrs...))
}

// RecordTopicOperation records a topic operation metric
func (m *OTelMetrics) RecordTopicOperation(ctx context.Context, operation string, success bool) {
	attrs := []attribute.KeyValue{
		attribute.String("operation", operation),
		attribute.Bool("success", success),
	}
	m.TopicOperationsTotal.Add(ctx, 1, metric.WithAttributes(attrs...))
}

// RecordNewsOperation records a news operation metric
func (m *OTelMetrics) RecordNewsOperation(ctx context.Context, operation string, success bool) {
	attrs := []attribute.KeyValue{
		attribute.String("operation", operation),
		attribute.Bool("success", success),
	}
	m.NewsOperationsTotal.Add(ctx, 1, metric.WithAttributes(attrs...))
}
