package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	mathRand "math/rand" // Alias für math/rand, um Konflikte mit crypto/rand zu vermeiden
	"os"
	"path/filepath"
	"time"

	dd_logrus "github.com/DataDog/dd-trace-go/contrib/sirupsen/logrus/v2"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	logtel "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

func generateInstanceID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatalf("Failed to generate instance ID: %v", err)
	}
	return hex.EncodeToString(b)
}

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

	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().IsValid() {
		ctx, span = otel.Tracer("otel-go-example").Start(ctx, "new-span")
		defer span.End()
	}
	traceID := span.SpanContext().TraceID().String()
	spanID := span.SpanContext().SpanID().String()

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("emitter-go"),
		semconv.ServiceVersionKey.String("1.0.0"),
		semconv.ServiceInstanceIDKey.String(generateInstanceID()),
		attribute.String("dd.trace_id", traceID),
		attribute.String("dd.span_id", spanID),
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

func setupFileLogger() (*os.File, error) {
	logDir := "logs"
	logFile := filepath.Join(logDir, "app.log")

	// Ensure the logs directory exists
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	// Open the log file
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	// Set logrus output to the file
	logrus.SetOutput(file)

	// Return the file for deferred closing
	return file, nil
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

func generateRandomLog(ctx context.Context) {
	logMessages := []string{
		"User logged in",
		"File uploaded",
		"Database updated",
		"Cache cleared",
		"Service started",
	}
	message := logMessages[mathRand.Intn(len(logMessages))]
	span := trace.SpanFromContext(ctx)
	traceID := span.SpanContext().TraceID().String()
	spanID := span.SpanContext().SpanID().String()

	logrus.WithContext(ctx).WithFields(logrus.Fields{
		"dd.trace_id": traceID,
		"dd.span_id":  spanID,
		"env":         "production",
		"service":     "emitter-go",
		"version":     "1.0.0",
	}).Info(message)
	log.Println(message)
}

func newResource() (*resource.Resource, error) {
	span := trace.SpanFromContext(context.Background())
	traceID := span.SpanContext().TraceID().String()
	spanID := span.SpanContext().SpanID().String()

	return resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceName("emitter-go"),
			semconv.ServiceVersion("0.1.0"),
			attribute.String("dd.trace_id", traceID),
			attribute.String("dd.span_id", spanID),
		))
}

func newLoggerProvider(ctx context.Context, res *resource.Resource) (*logtel.LoggerProvider, error) {
	exporter, err := otlploggrpc.New(ctx)
	if err != nil {
		return nil, err
	}
	processor := logtel.NewSimpleProcessor(exporter)
	provider := logtel.NewLoggerProvider(
		logtel.WithResource(res),
		logtel.WithProcessor(processor),
	)
	return provider, nil
}

func main() {
	ctx := context.Background()
	res, err := newResource()

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.AddHook(&dd_logrus.DDContextLogHook{})

	// Setup file logger
	logFile, err := setupFileLogger()
	if err != nil {
		log.Fatalf("Failed to setup file logger: %v", err)
	}
	defer logFile.Close()

	// Redirect standard log output to the same file
	log.SetOutput(logFile)

	if err != nil {
		log.Fatalf("Failed to create resource: %v", err)
	}
	loggerProvider, err := newLoggerProvider(ctx, res)
	if err != nil {
		log.Fatalf("Failed to create logger provider: %v", err)
	}
	tracerProvider, meterProvider, err := setupOpenTelemetry(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize OpenTelemetry: %v", err)
	}

	defer func() {
		if tracerProvider != nil {
			if err := tracerProvider.Shutdown(ctx); err != nil {
				log.Printf("Error shutting down TracerProvider: %v", err)
			}
		}
		if meterProvider != nil {
			if err := meterProvider.Shutdown(ctx); err != nil {
				log.Printf("Error shutting down MeterProvider: %v", err)
			}
		}
		if loggerProvider != nil {
			if err := loggerProvider.Shutdown(ctx); err != nil {
				fmt.Println(err)
			}
		}
	}()

	tracer := otel.Tracer("otel-go-example")
	meter := otel.Meter("otel-go-example")

	logrus.WithContext(ctx).Info("Go logs and traces connected!")

	counter, err := meter.Int64Counter("example_counter")
	if err != nil {
		log.Fatalf("Failed to create counter: %v", err)
	}
	histogram, err := meter.Float64Histogram("example_histogram")
	if err != nil {
		log.Fatalf("Failed to create histogram: %v", err)
	}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		ctx, span := tracer.Start(ctx, "random-trace")
		generateRandomTrace(tracer)
		generateRandomLog(ctx)
		span.End()

		counter.Add(ctx, 1)
		histogram.Record(ctx, float64(mathRand.Intn(100)))
	}
}
