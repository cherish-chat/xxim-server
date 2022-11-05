package xtrace

import (
	"context"
	"github.com/zeromicro/go-zero/core/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type StartFuncSpanOpt struct {
	carrier  propagation.TextMapCarrier
	spanKind *oteltrace.SpanKind
}
type StartFuncSpanOptFunc func(opt *StartFuncSpanOpt)

func StartFuncSpanWithCarrier(carrier propagation.TextMapCarrier) StartFuncSpanOptFunc {
	return func(opt *StartFuncSpanOpt) {
		opt.carrier = carrier
	}
}
func StartFuncSpanWithKind(spanKind oteltrace.SpanKind) StartFuncSpanOptFunc {
	return func(opt *StartFuncSpanOpt) {
		opt.spanKind = &spanKind
	}
}
func StartFuncSpan(ctx context.Context, name string, f func(context.Context), opts ...StartFuncSpanOptFunc) {
	var opt StartFuncSpanOpt
	for _, o := range opts {
		o(&opt)
	}
	var kv []attribute.KeyValue
	if opt.carrier != nil {
		for _, k := range opt.carrier.Keys() {
			v := opt.carrier.Get(k)
			kv = append(kv, attribute.String(k, v))
		}
	}
	var spanKind oteltrace.SpanKind
	if opt.spanKind != nil {
		spanKind = *opt.spanKind
	} else {
		spanKind = oteltrace.SpanKindInternal
	}
	tracer := otel.GetTracerProvider().Tracer(trace.TraceName)
	spanCtx, span := tracer.Start(ctx, name,
		oteltrace.WithSpanKind(spanKind),
		oteltrace.WithAttributes(kv...),
	)
	defer span.End()
	f(spanCtx)
}

func RunWithTrace(
	traceId string,
	spanName string,
	f func(ctx context.Context),
	carrier propagation.TextMapCarrier,
) {
	var attributes []attribute.KeyValue
	attributes = append(attributes, attribute.String("traceId", traceId))
	if carrier != nil {
		for _, k := range carrier.Keys() {
			v := carrier.Get(k)
			attributes = append(attributes, attribute.String(k, v))
		}
	}
	tracer := otel.GetTracerProvider().Tracer(trace.TraceName)
	traceIDFromHex, _ := oteltrace.TraceIDFromHex(traceId)
	ctx := oteltrace.ContextWithSpanContext(context.Background(), oteltrace.NewSpanContext(oteltrace.SpanContextConfig{
		TraceID: traceIDFromHex,
	}))
	spanCtx, span := tracer.Start(
		ctx,
		spanName,
		oteltrace.WithAttributes(attributes...),
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
