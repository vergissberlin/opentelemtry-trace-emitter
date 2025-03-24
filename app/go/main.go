package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	mathRand "math/rand" // Alias für math/rand, um Konflikte mit crypto/rand zu vermeiden
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

// Generiert eine zufällige ID für die Instanz
func generateInstanceID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatalf("Failed to generate instance ID: %v", err)
	}
	return hex.EncodeToString(b)
}

// Initialisiert OpenTelemetry für Tracing und Metriken
func setupOpenTelemetry(ctx context.Context) (*sdktrace.TracerProvider, *sdkmetric.MeterProvider, error) {
	collectorURL := os.Getenv("OTEL_COLLECTOR_ENDPOINT")
	if collectorURL == "" {
		collectorURL = "http://localhost:4317"
	}

	println("Using OpenTelemetry Collector endpoint:", collectorURL)

	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	metricExporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithInsecure())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create metric exporter: %w", err)
	}

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("emitter-go"),
		semconv.ServiceVersionKey.String("1.0.0"),
		semconv.ServiceInstanceIDKey.String(generateInstanceID()),
	)

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tracerProvider)

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter)),
		sdkmetric.WithResource(res),
	)
	otel.SetMeterProvider(meterProvider)

	return tracerProvider, meterProvider, nil
}

func generateRandomTrace(tracer trace.Tracer) {
	ctx, span := tracer.Start(context.Background(), "random-trace")
	defer span.End()

	operations := []string{"database-query", "api-call", "cache-fetch", "file-read"}
	operation := operations[mathRand.Intn(len(operations))]
	_, subSpan := tracer.Start(ctx, operation)
	time.Sleep(time.Duration(mathRand.Intn(500)) * time.Millisecond)
	subSpan.End()

	fmt.Printf("Generated trace: %s\n", operation)
}

func generateRandomLog() {
	logMessages := []string{
		"User logged in",
		"File uploaded",
		"Database updated",
		"Cache cleared",
		"Service started",
	}
	message := logMessages[mathRand.Intn(len(logMessages))]
	log.Println(message)
}

func newResource() (*resource.Resource, error) {
	return resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceName("my-service"),
			semconv.ServiceVersion("0.1.0"),
		))
}

func newLoggerProvider(ctx context.Context, res *resource.Resource) (*log.LoggerProvider, error) {
	exporter, err := otlploggrpc.New(ctx)
	if err != nil {
		return nil, err
	}
	processor := log.NewBatchProcessor(exporter)
	provider := log.NewLoggerProvider(
		log.WithResource(res),
		log.WithProcessor(processor),
	)
	return provider, nil
}

func main() {
	ctx := context.Background()
	loggerProvider, err := newLoggerProvider(ctx, res)
	tracerProvider, meterProvider, err := setupOpenTelemetry(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize OpenTelemetry: %v", err)
	}

	defer func() {
		if err := tracerProvider.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down TracerProvider: %v", err)
		}
		if err := meterProvider.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down MeterProvider: %v", err)
		}
		if err := loggerProvider.Shutdown(ctx); err != nil {
			fmt.Println(err)
		}
	}()

	tracer := otel.Tracer("otel-go-example")

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		generateRandomTrace(tracer)
		generateRandomLog()
	}
}
