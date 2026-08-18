package main

import (
	"context"
	"log"

	"github.com/tuanta7/ciam/internal/config"
	"github.com/tuanta7/ciam/pkg/o11y"
)

func initMonitor(ctx context.Context, cfg *config.EnvConfig) {
	if !cfg.EnableMetrics {
		o11y.InitNoopMeterProvider()
	} else {
		_, err := o11y.InitMeterProvider(ctx, cfg.ServiceName, nil)
		if err != nil {
			log.Fatalf("Failed to initialize meter provider: %v", err)
		}
	}

	if !cfg.EnableTracing {
		o11y.InitNoopTracerProvider()
	} else {
		_, err := o11y.InitTracerProvider(ctx, cfg.ServiceName, nil)
		if err != nil {
			log.Fatalf("Failed to initialize tracer provider: %v", err)
		}
	}
}
