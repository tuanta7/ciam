package main

import (
	"context"
	"log"

	"github.com/tuanta7/ciam/internal/config"
	"github.com/tuanta7/ciam/pkg/otelx"
)

func initMonitor(ctx context.Context, cfg *config.EnvConfig) {
	if !cfg.EnableMetrics {
		otelx.InitNoopMeterProvider()
	} else {
		_, err := otelx.InitMeterProvider(ctx, cfg.ServiceName, nil)
		if err != nil {
			log.Fatalf("Failed to initialize meter provider: %v", err)
		}
	}

	if !cfg.EnableTracing {
		otelx.InitNoopTracerProvider()
	} else {
		_, err := otelx.InitTracerProvider(ctx, cfg.ServiceName, nil)
		if err != nil {
			log.Fatalf("Failed to initialize tracer provider: %v", err)
		}
	}
}
