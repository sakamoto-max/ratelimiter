package opentel

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

type traceKey string

var TraceKey traceKey = "TRACE_KEY"

func InitTracing(ctx context.Context) (*sdktrace.TracerProvider, error) {
	// exporter, err := otlptracehttp.New(
	// 	ctx,
	// 	otlptracehttp.WithEndpoint("jaeger:4318"),
	// )
	// if err != nil {
	// 	return nil, err
	// }

	// // res, err := resource.Merge(
	// // 	resource.Default(),
	// // 	resource.NewWithAttributes(
	// // 		// semconv.SchemaURL,
	// // 		semconv.ServiceName("rate-limiter"),
	// // 	),
	// // )

	res, err := resource.New(
		ctx,
		resource.WithAttributes(
			attribute.String("service.name", "rate-limiter"),
		),
	)
	if err != nil {
		return nil, err
	}

	// tp := sdktrace.NewTracerProvider(
	// 	sdktrace.WithBatcher(exporter),
	// 	sdktrace.WithResource(res),
	// )

	// otel.SetTracerProvider(tp)
	// otel.SetTextMapPropagator(
	// 	propagation.NewCompositeTextMapPropagator(
	// 		propagation.TraceContext{},
	// 		propagation.Baggage{},
	// 	),
	// )

	exporter, err := otlptracehttp.New(
		ctx,
		otlptracehttp.WithEndpoint("jaeger:4318"),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("create exporter: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		// sdktrace.WithSyncer(exporter),
		//  // Debugging only
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)

	return tp, nil
}

func CreateNewSpan(ctx context.Context, serviceName string) (context.Context, trace.Span) {
	tracerFromCtx := ctx.Value(TraceKey)

	var tracer trace.Tracer
	if tracerFromCtx == nil {
		tracer = otel.Tracer("RATE_LIMITER")
		ctx = context.WithValue(ctx, TraceKey, tracer)
	} else {
		tracer = tracerFromCtx.(trace.Tracer)
	}

	ctx, span := tracer.Start(ctx, serviceName)

	return context.WithValue(ctx, TraceKey, tracer), span
}
