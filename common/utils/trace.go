package utils

import (
	"context"
	"github.com/zeromicro/go-zero/core/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type xTrace struct {
}

var Trace = &xTrace{}

func (x *xTrace) Span(ctx context.Context, spanName string, do func(context.Context) error, attributes map[string]string) {
	if do == nil {
		return
	}
	if attributes == nil {
		attributes = make(map[string]string)
	}
	var kvs []attribute.KeyValue
	for k, v := range attributes {
		kvs = append(kvs, attribute.String(k, v))
	}
	spanKind := oteltrace.SpanKindInternal
	tracer := otel.GetTracerProvider().Tracer(trace.TraceName)
	spanCtx, span := tracer.Start(ctx, spanName,
		oteltrace.WithSpanKind(spanKind),
		oteltrace.WithAttributes(kvs...),
	)
	defer span.End()
	err := do(spanCtx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
	} else {
		span.SetStatus(codes.Ok, "")
	}
}
