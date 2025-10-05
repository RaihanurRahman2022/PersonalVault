package monitoring

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type Metrics struct {
	//Application Metrics
	requestCounter    metric.Int64Counter
	requestDuration   metric.Float64Histogram
	activeConnections metric.Int64UpDownCounter
	errorCounter      metric.Int64Counter
	dbConnections     metric.Int64UpDownCounter
	fileOperations    metric.Int64Counter

	// System Metrics
	systemMetrics *SystemMetrics
}

func NewMetrics() *Metrics {
	meter := otel.Meter("personal-vault-server")

	//Application Metrics
	requestCounter, _ := meter.Int64Counter("http_requests_total",
		metric.WithDescription("Total number of HTTP requests"))

	requestDuration, _ := meter.Float64Histogram(
		"http_request_duration_seconds",
		metric.WithDescription("HTTP request duration in seconds"),
	)

	activeConnections, _ := meter.Int64UpDownCounter(
		"http_active_connections",
		metric.WithDescription("Number of active HTTP connections"),
	)

	errorCounter, _ := meter.Int64Counter(
		"http_errors_total",
		metric.WithDescription("Total number of HTTP errors"),
	)

	dbConnections, _ := meter.Int64UpDownCounter(
		"db_connections_active",
		metric.WithDescription("Number of active database connections"),
	)

	fileOperations, _ := meter.Int64Counter(
		"file_operations_total",
		metric.WithDescription("Total number of file operations"),
	)

	// Initialize system metrics
	systemMetrics := NewSystemMetrics()

	return &Metrics{
		requestCounter:    requestCounter,
		requestDuration:   requestDuration,
		activeConnections: activeConnections,
		errorCounter:      errorCounter,
		dbConnections:     dbConnections,
		fileOperations:    fileOperations,
		systemMetrics:     systemMetrics,
	}
}

func (m *Metrics) RecordRequest(method, path string, statusCode int, duration time.Duration) {
	ctx := context.Background()

	// Record request count
	m.requestCounter.Add(ctx, 1,
		metric.WithAttributes(
			attribute.String("method", method),
			attribute.String("path", path),
			attribute.Int("status_code", statusCode),
		),
	)

	// Record request duration
	m.requestDuration.Record(ctx, duration.Seconds(),
		metric.WithAttributes(
			attribute.String("method", method),
			attribute.String("path", path),
		),
	)

	// Record errors
	if statusCode >= 400 {
		m.errorCounter.Add(ctx, 1,
			metric.WithAttributes(
				attribute.String("method", method),
				attribute.String("path", path),
				attribute.Int("status_code", statusCode),
			),
		)
	}
}

func (m *Metrics) IncrementActiveConnections() {
	m.activeConnections.Add(context.Background(), 1)
}

func (m *Metrics) DecrementActiveConnections() {
	m.activeConnections.Add(context.Background(), -1)
}

func (m *Metrics) IncrementDBConnections() {
	m.dbConnections.Add(context.Background(), 1)
}

func (m *Metrics) DecrementDBConnections() {
	m.dbConnections.Add(context.Background(), -1)
}

func (m *Metrics) RecordFileOperation(operation, fileType string) {
	m.fileOperations.Add(context.Background(), 1,
		metric.WithAttributes(
			attribute.String("operation", operation),
			attribute.String("file_type", fileType),
		),
	)
}

func (m *Metrics) StartSystemMonitoring() {
	m.systemMetrics.StartSystemMonitoring()
}
