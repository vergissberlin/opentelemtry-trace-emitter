package main

import (
	"context"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()

	// OTel gRPC Exporter
	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint("localhost:4317"),
		otlptracegrpc.WithDialOption(grpc.WithBlock()),
	)
	if err != nil {
		log.Fatalf("Failed to create exporter: %v", err)
	}
	defer func() {
		if err := exporter.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down exporter: %v", err)
		}
	}()

	// Tracer Provider
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewSchemaless(
			attribute.String("service.name", "go-log-service"),
		)),
	)
	otel.SetTracerProvider(tracerProvider)
	defer func() {
		if err := tracerProvider.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	// Create a Tracer
	tracer := otel.Tracer("example-logger")
	ctx, span := tracer.Start(ctx, "log-event")
	log.Println("Sending log event to OTel collector")
	span.AddEvent("Log message", trace.WithAttributes(attribute.String("message", "Hello OpenTelemetry!")))
	time.Sleep(1 * time.Second)
	span.End()
}
