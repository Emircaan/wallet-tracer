package opentelemetry

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.uber.org/zap"
)

func NewExporter(ctx context.Context) sdktrace.SpanExporter {
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://localhost:14268/api/traces")))

	if err != nil {
		zap.L().Fatal("failed to create exporter", zap.Error(err))
	}

	return exporter

}

func NewTraceProvider(exp sdktrace.SpanExporter) *sdktrace.TracerProvider {
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("account-service"),
		),
	)
	if err != nil {
		panic(err)
	}
	tp := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exp), sdktrace.WithResource(r))
	otel.SetTracerProvider(tp)

	return tp

}

func InitTelemetry(ctx context.Context) func() {
	exp := NewExporter(ctx)
	tp := NewTraceProvider(exp)
	return func() {
		if err := tp.Shutdown(ctx); err != nil {
			zap.L().Fatal("failed to shutdown provider", zap.Error(err))
		}
	}

}
