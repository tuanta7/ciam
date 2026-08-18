package o11y

import (
	"log/slog"

	"go.opentelemetry.io/contrib/bridges/otelslog"
)

type Logger struct {
	*slog.Logger
}

func NewLogger(name string) *Logger {
	return &Logger{
		Logger: otelslog.NewLogger(name),
	}
}
