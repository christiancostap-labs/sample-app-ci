package main

import (
	"context"
	"go.uber.org/zap"

	"go.opentelemetry.io/otel/trace"
)

func loggerWithCorrelateTrace(ctx context.Context) *zap.Logger {
	sc := trace.SpanFromContext(ctx).SpanContext()

	if !sc.IsValid() {
		return zap.L()
	}

	return zap.L().With(
		zap.String("trace_id", sc.TraceID().String()),
		zap.String("span_id", sc.SpanID().String()),
	)
}
