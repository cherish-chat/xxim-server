package xtrace

import (
	"context"
	"github.com/zeromicro/go-zero/core/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
)

func StartFuncSpan(ctx context.Context, name string, f func(context.Context), kv ...attribute.KeyValue) {
	tracer := otel.GetTracerProvider().Tracer(trace.TraceName)
	spanCtx, span := tracer.Start(ctx, name,
		oteltrace.WithSpanKind(oteltrace.SpanKindInternal),
		oteltrace.WithAttributes(kv...),
	)
	defer span.End()
	f(spanCtx)
}

func RunWithTrace(
	traceId string,
	spanName string,
	f func(ctx context.Context),
) {
	tracer := otel.GetTracerProvider().Tracer(trace.TraceName)
	traceIDFromHex, _ := oteltrace.TraceIDFromHex(traceId)
	ctx := oteltrace.ContextWithSpanContext(context.Background(), oteltrace.NewSpanContext(oteltrace.SpanContextConfig{
		TraceID: traceIDFromHex,
	}))
	spanCtx, span := tracer.Start(
		ctx,
		spanName,
	)
	defer span.End()
	f(spanCtx)
}

func TraceIdFromContext(ctx context.Context) string {
	spanCtx := oteltrace.SpanContextFromContext(ctx)
	if spanCtx.HasTraceID() {
		traceId := spanCtx.TraceID().String()
		return traceId
	}
	return ""
}
