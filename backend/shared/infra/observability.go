package infra

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

// InitMetrics initializes Prometheus metrics exporter
func InitMetrics(ctx context.Context, serviceName string) (metric.MeterProvider, error) {
	exporter, err := prometheus.New()
	if err != nil {
		return nil, err
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
		),
	)
	if err != nil {
		return nil, err
	}

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(exporter),
	)

	otel.SetMeterProvider(meterProvider)
	return meterProvider, nil
}

// MetricsService wraps metrics functionality
type MetricsService struct {
	Meter metric.Meter
}

// NewMetricsService creates a new metrics service
func NewMetricsService(serviceName string) *MetricsService {
	meter := otel.Meter(serviceName)
	return &MetricsService{Meter: meter}
}
