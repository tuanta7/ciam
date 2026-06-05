package monitor

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.39.0"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
	"google.golang.org/grpc"
)

func InitNoopTracerProvider() trace.TracerProvider {
	provider := noop.NewTracerProvider()
	otel.SetTracerProvider(provider)
	return provider
}

func InitTracerProvider(ctx context.Context, serviceName string, gc *grpc.ClientConn) (*sdktrace.TracerProvider, error) {
	res, err := resource.New(ctx,
		resource.WithAttributes(semconv.ServiceNameKey.String(serviceName)),
	)
	if err != nil {
		return nil, err
	}

	otlpExporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithGRPCConn(gc),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(otlpExporter, sdktrace.WithBatchTimeout(5*time.Second)),
	)

	otel.SetTracerProvider(tracerProvider)
	return tracerProvider, nil
}

// InitPropagator sets the global propagator used for cross-service context propagation.
// Call this once during service startup (before handling requests / making outbound calls).
func InitPropagator() {
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)
}
