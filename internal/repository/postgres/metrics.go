package postgres

import "go.opentelemetry.io/otel/metric"

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
