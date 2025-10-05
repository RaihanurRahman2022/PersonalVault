package monitoring

import (
	"context"
	"fmt"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.32.0"
)

type TelemetryConfig struct {
	ServiceName    string
	ServiceVersion string
	Environment    string
	JaegerEndpoint string
}

func InitTelemetry(config TelemetryConfig) (*trace.TracerProvider, *metric.MeterProvider, error) {
	// Create resource
	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(config.ServiceName),
			semconv.ServiceVersionKey.String(config.ServiceVersion),
			semconv.DeploymentEnvironmentKey.String(config.Environment),
		),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Initialize tracer provider
	tp, err := initTracerProvider(res, config.JaegerEndpoint)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize tracer provider: %w", err)
	}

	// Initialize meter provider
	mp, err := initMeterProvider(res)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize meter provider: %w", err)
	}

	// Set global providers
	otel.SetTracerProvider(tp)
	otel.SetMeterProvider(mp)

	return tp, mp, nil
}

func initTracerProvider(res *resource.Resource, jaegerEndpoint string) (*trace.TracerProvider, error) {
	// Create Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerEndpoint)))
	if err != nil {
		return nil, fmt.Errorf("failed to create Jaeger exporter: %w", err)
	}

	// Create tracer provider
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(res),
	)

	return tp, nil
}

func initMeterProvider(res *resource.Resource) (*metric.MeterProvider, error) {
	// Create Prometheus exporter
	exp, err := prometheus.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create Prometheus exporter: %w", err)
	}

	// Create meter provider
	mp := metric.NewMeterProvider(
		metric.WithReader(exp),
		metric.WithResource(res),
	)

	return mp, nil
}

func ShutdownTelemetry(tp *trace.TracerProvider, mp *metric.MeterProvider) {
	if tp != nil {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}
	if mp != nil {
		if err := mp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down meter provider: %v", err)
		}
	}
}
