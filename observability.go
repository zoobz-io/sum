package sum

import (
	"context"
	"fmt"

	"github.com/zoobz-io/aperture"
	"github.com/zoobz-io/capitan"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// WithObservability configures OTEL exporters and an aperture bridge for the service.
// name is the OTEL service name. endpoint is the OTLP collector endpoint
// (e.g., "http://localhost:4318").
// Aperture defaults to logging all capitan events. Use Observability()
// to access the aperture instance for custom schema configuration via Apply.
func (s *Service) WithObservability(ctx context.Context, name, endpoint string) error {
	res, err := resource.New(ctx,
		resource.WithAttributes(semconv.ServiceName(name)),
	)
	if err != nil {
		return fmt.Errorf("observability resource: %w", err)
	}

	traceExp, err := otlptracehttp.New(ctx, otlptracehttp.WithEndpointURL(endpoint))
	if err != nil {
		return fmt.Errorf("observability trace exporter: %w", err)
	}

	metricExp, err := otlpmetrichttp.New(ctx, otlpmetrichttp.WithEndpointURL(endpoint))
	if err != nil {
		return fmt.Errorf("observability metric exporter: %w", err)
	}

	logExp, err := otlploghttp.New(ctx, otlploghttp.WithEndpointURL(endpoint))
	if err != nil {
		return fmt.Errorf("observability log exporter: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExp),
		sdktrace.WithResource(res),
	)
	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExp)),
		sdkmetric.WithResource(res),
	)
	lp := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(sdklog.NewBatchProcessor(logExp)),
		sdklog.WithResource(res),
	)

	a, err := aperture.New(capitan.Default(), lp, mp, tp)
	if err != nil {
		// Shutdown providers on failure.
		tp.Shutdown(ctx)
		mp.Shutdown(ctx)
		lp.Shutdown(ctx)
		return fmt.Errorf("observability aperture: %w", err)
	}

	s.mu.Lock()
	s.aperture = a
	s.providers = &otelProviders{
		trace:  tp,
		meter:  mp,
		logger: lp,
	}
	s.mu.Unlock()

	return nil
}

// Observability returns the aperture instance for custom schema configuration.
// Returns nil if WithObservability has not been called.
func (s *Service) Observability() *aperture.Aperture {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.aperture
}

// otelProviders holds OTEL SDK providers for lifecycle management.
type otelProviders struct {
	trace  *sdktrace.TracerProvider
	meter  *sdkmetric.MeterProvider
	logger *sdklog.LoggerProvider
}

// Shutdown flushes and shuts down all OTEL providers in order: trace, meter, log.
func (p *otelProviders) Shutdown(ctx context.Context) error {
	var firstErr error
	if err := p.trace.Shutdown(ctx); err != nil && firstErr == nil {
		firstErr = err
	}
	if err := p.meter.Shutdown(ctx); err != nil && firstErr == nil {
		firstErr = err
	}
	if err := p.logger.Shutdown(ctx); err != nil && firstErr == nil {
		firstErr = err
	}
	return firstErr
}
