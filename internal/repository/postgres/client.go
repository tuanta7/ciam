package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.39.0"
	"go.opentelemetry.io/otel/trace"
)

// InstrumentedPool is a wrapper around pgxpool.Pool that instruments the pool.
// It implements the store.DBTX interface.
type InstrumentedPool struct {
	pool   *pgxpool.Pool
	tracer trace.Tracer
	meter  metric.Meter
}

func NewInstrumentedPool(ctx context.Context, dsn string) (*InstrumentedPool, error) {
	pgxPool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	if err := pgxPool.Ping(ctx); err != nil {
		return nil, err
	}

	p := &InstrumentedPool{
		pool:   pgxPool,
		tracer: otel.Tracer("postgres_tracer"),
		meter:  otel.Meter("postgres_meter"),
	}

	err = initMetrics(p.meter)
	return p, err
}

func (p *InstrumentedPool) Close() {
	p.pool.Close()
}

func (p *InstrumentedPool) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	ctx, span := p.tracer.Start(ctx, "postgres_exec", trace.WithAttributes(
		semconv.DBSystemNamePostgreSQL,
		semconv.DBQueryText(sql),
	))
	defer span.End()

	commandTag, err := p.pool.Exec(ctx, sql, arguments...)
	if err != nil {
		span.RecordError(err)
		return commandTag, err
	}

	return commandTag, nil
}

func (p *InstrumentedPool) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	start := time.Now()
	ctx, span := p.tracer.Start(ctx, "postgres_query", trace.WithAttributes(
		semconv.DBSystemNamePostgreSQL,
		semconv.DBQueryText(sql),
	))
	defer span.End()

	rows, err := p.pool.Query(ctx, sql, args...)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	queryDuration.Record(ctx, time.Since(start).Seconds(), metric.WithAttributes(
		attribute.String("db.operation", "query"),
	))

	return rows, err
}

func (p *InstrumentedPool) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	start := time.Now()
	ctx, span := p.tracer.Start(ctx, "postgres_query_row", trace.WithAttributes(
		semconv.DBSystemNamePostgreSQL,
		semconv.DBQueryText(sql),
	))
	defer span.End()

	row := p.pool.QueryRow(ctx, sql, args...)

	queryDuration.Record(ctx, time.Since(start).Seconds(), metric.WithAttributes(
		attribute.String("db.operation", "query_row"),
	))

	return row
}
