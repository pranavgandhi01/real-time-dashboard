package observability

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	v1170 "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func InitTracing(serviceName, jaegerEndpoint string) (*trace.TracerProvider, error) {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerEndpoint)))
	if err != nil {
		return nil, err
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(resource.NewWithAttributes(
			v1170.SchemaURL,
			v1170.ServiceName(serviceName),
		)),
	)

	otel.SetTracerProvider(tp)
	return tp, nil
}