package monitor

import (
	"context"
	"os"

	"go.opentelemetry.io/contrib/bridges/otelzap"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.39.0"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

type Logger struct {
	*zap.Logger
}

func NewNoopLogger() *Logger {
	return &Logger{
		Logger: zap.NewNop(),
	}
}

func NewLogger(ctx context.Context, serviceName string, gc *grpc.ClientConn, level ...zapcore.Level) (*Logger, error) {
	level = append(level, zapcore.InfoLevel)

	cfg := zap.NewProductionConfig()
	cfg.Encoding = "json"
	cfg.Level = zap.NewAtomicLevelAt(level[0])

	var provider *sdklog.LoggerProvider
	var err error

	if isOTLPLogger() {
		provider, err = initLoggerProvider(ctx, serviceName, gc)
		if err != nil {
			return nil, err
		}
	}

	zl, err := cfg.Build(
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			if isOTLPLogger() {
				return otelzap.NewCore("", otelzap.WithLoggerProvider(provider))
			}
			return core
		}),
	)
	if err != nil {
		return nil, err
	}

	return &Logger{
		Logger: zl,
	}, nil
}

func isOTLPLogger() bool {
	return os.Getenv("LOGGER_OUTPUT") == "otlp"
}

func initLoggerProvider(ctx context.Context, serviceName string, gc *grpc.ClientConn) (*sdklog.LoggerProvider, error) {
	res, err := resource.New(ctx,
		resource.WithAttributes(semconv.ServiceNameKey.String(serviceName)),
	)
	if err != nil {
		return nil, err
	}

	exporter, err := otlploggrpc.New(ctx, otlploggrpc.WithGRPCConn(gc))
	if err != nil {
		return nil, err
	}

	return sdklog.NewLoggerProvider(
		sdklog.WithResource(res),
		sdklog.WithProcessor(sdklog.NewBatchProcessor(exporter)),
	), nil
}

func (l *Logger) Close() error {
	return l.Logger.Sync()
}
