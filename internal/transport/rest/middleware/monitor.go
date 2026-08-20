package middleware

import (
	"log/slog"
	"net/http"

	"github.com/tuanta7/ciam/pkg/otelx"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.39.0"
	"go.opentelemetry.io/otel/trace"
)

var (
	requestTotal metric.Int64Counter
	requestError metric.Int64Counter
)

func InitMetricsMiddleware(meter metric.Meter) (err error) {
	requestTotal, err = meter.Int64Counter("http.server.requests",
		metric.WithDescription("Number of HTTP requests"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return err
	}

	requestError, err = meter.Int64Counter("http.server.errors",
		metric.WithDescription("Number of HTTP errors"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return err
	}

	return nil
}

type responseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (w *responseWriter) WriteHeader(code int) {
	w.StatusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func WithMetric(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := newResponseWriter(w)
		next.ServeHTTP(rw, r)

		requestTotal.Add(r.Context(), 1, metric.WithAttributes(
			attribute.String("method", r.Method),
			attribute.String("path", r.URL.Path),
			attribute.Int("status_code", rw.StatusCode),
		))

		if rw.StatusCode >= 400 {
			requestError.Add(r.Context(), 1, metric.WithAttributes(
				attribute.String("method", r.Method),
				attribute.String("path", r.URL.Path),
				attribute.Int("status_code", rw.StatusCode),
			))
		}
	})
}

func WithTrace(tracer trace.Tracer, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Continue trace from upstream by extracting the remote parent span context.
		parentCtx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))
		ctx, span := tracer.Start(parentCtx, "http.request")
		defer span.End()

		rw := newResponseWriter(w)
		next.ServeHTTP(rw, r.WithContext(ctx))

		span.SetAttributes(
			semconv.HTTPResponseStatusCode(rw.StatusCode),
			semconv.HTTPRequestMethodOriginal(r.Method),
			semconv.URLPath(r.URL.Path),
		)
	})
}

func WithLog(logger *otelx.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Request received",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
		)

		next.ServeHTTP(w, r)
	})
}

func WithTelemetry(tracer trace.Tracer, logger *otelx.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		WithLog(logger, WithTrace(tracer, WithMetric(next))).ServeHTTP(w, r)
	})
}
