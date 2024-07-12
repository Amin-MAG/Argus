package tracing

import (
	"context"
	"go.opentelemetry.io/otel/propagation"

	propagators "go.opentelemetry.io/contrib/propagators/jaeger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"

	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

var globalTracerName string

var tp *trace.TracerProvider

// TracerProvider Creates the Jaeger exporter
// TracerProvider returns an OpenTelemetry TracerProvider configured to use
// the Jaeger exporter that will send spans to the provided url. The returned
// TracerProvider will also use a Resource configured with all the information
// about the application.
func TracerProvider(serviceName, serviceVersion string, samplerRatio float64, endpoint string) (*trace.TracerProvider, error) {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(endpoint)))
	if err != nil {
		return nil, err
	}

	return trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithSampler(trace.ParentBased(trace.TraceIDRatioBased(samplerRatio))),
		// Record information about this application in a Resource.
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersionKey.String(serviceVersion),
		)),
	), nil
}

// InitTracing Sets globalTracerName and provide the tracer, then sets Jaeger Propagation
func InitTracing(serviceName, serviceVersion string, samplerRatio float64, endpoint string) error {
	// Set the global tracer name
	globalTracerName = serviceName

	// Create the tracer
	var err error
	tp, err = TracerProvider(serviceName, serviceVersion, samplerRatio, endpoint)
	if err != nil {
		return err
		//logger.FatalS("could not create tracer provider", "error", err.Error())
	}

	// Register our TracerProvider as the global so any imported
	// instrumentation in the future will default to using it.
	otel.SetTracerProvider(tp)

	// Config the propagation
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
		propagators.Jaeger{},
	))

	return nil
}

func Shutdown(ctx context.Context) error {
	return tp.Shutdown(ctx)
}

// TracerName Returns the global trace name
func TracerName() string {
	return globalTracerName
}
