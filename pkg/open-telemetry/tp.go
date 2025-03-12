package opentelemetry

import (
	"context"

	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func NewExporter(ctx context.Context) {

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

	return sdktrace.NewTracerProvider(sdktrace.WithBatcher(exp), sdktrace.WithResource(r))

}
