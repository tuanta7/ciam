package postgres

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.39.0"
	"go.opentelemetry.io/otel/trace"
)

// InstrumentedPool wraps *sql.DB to instrument queries.
// It implements boil.ContextExecutor for sqlboiler.
type InstrumentedPool struct {
	*sql.DB
	tracer trace.Tracer
	meter  metric.Meter
}

func NewInstrumentedPool(ctx context.Context, dsn string) (*InstrumentedPool, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	p := &InstrumentedPool{
		DB:     db,
		tracer: otel.Tracer("postgres_tracer"),
		meter:  otel.Meter("postgres_meter"),
	}

	err = initMetrics(p.meter)
	return p, err
}

func (p *InstrumentedPool) Close() {
	p.DB.Close()
}

func (p *InstrumentedPool) ExecContext(ctx context.Context, sqlQuery string, args ...any) (sql.Result, error) {
	ctx, span := p.tracer.Start(ctx, "postgres_exec", trace.WithAttributes(
		semconv.DBSystemNamePostgreSQL,
		semconv.DBQueryText(sqlQuery),
	))
	defer span.End()

	result, err := p.DB.ExecContext(ctx, sqlQuery, args...)
	if err != nil {
		span.RecordError(err)
		return result, err
	}

	return result, nil
}

func (p *InstrumentedPool) QueryContext(ctx context.Context, sqlQuery string, args ...any) (*sql.Rows, error) {
	start := time.Now()
	ctx, span := p.tracer.Start(ctx, "postgres_query", trace.WithAttributes(
		semconv.DBSystemNamePostgreSQL,
		semconv.DBQueryText(sqlQuery),
	))
	defer span.End()

	rows, err := p.DB.QueryContext(ctx, sqlQuery, args...)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	queryDuration.Record(ctx, time.Since(start).Seconds(), metric.WithAttributes(
		attribute.String("db.operation", "query"),
	))

	return rows, err
}

func (p *InstrumentedPool) QueryRowContext(ctx context.Context, sqlQuery string, args ...any) *sql.Row {
	start := time.Now()
	ctx, span := p.tracer.Start(ctx, "postgres_query_row", trace.WithAttributes(
		semconv.DBSystemNamePostgreSQL,
		semconv.DBQueryText(sqlQuery),
	))
	defer span.End()

	row := p.DB.QueryRowContext(ctx, sqlQuery, args...)

	queryDuration.Record(ctx, time.Since(start).Seconds(), metric.WithAttributes(
		attribute.String("db.operation", "query_row"),
	))

	return row
}

var (
	queryDuration metric.Float64Histogram
)

func initMetrics(meter metric.Meter) error {
	var err error

	queryDuration, err = meter.Float64Histogram("db.client.operation.duration",
		metric.WithDescription("Duration of database operations"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return err
	}

	return nil
}
