# Emit traces with Go

This example shows how to emit traces with Go using the OpenTelemetry Go SDK. It uses the OpenTelemetry Collector as a receiver and exporter.

## Usage

### 1. Pull the image

```bash
docker pull ghcr.io/vergissberlin/opentelemetry-trace-emitter:go
```

### 2. Run the image

```bash
docker run -it --rm \
    -e OTEL_COLLECTOR_ENDPOINT=http://localhost:4317 \
    ghcr.io/vergissberlin/opentelemetry-trace-emitter:go
```

### 3. Environment Variables

|   Environment Variable    |                  Description                   |         Default         |
| ------------------------- | ---------------------------------------------- | ----------------------- |
| `OTEL_COLLECTOR_ENDPOINT` | The endpoint to send the spans to.             | `http://localhost:4317` |
| `TRACE_INTERVAL` | The timeout for the exporter.                  | `5000` |
| `SPAN_COUNT_MAX` | The maximum number of spans per trace to send. | `12` |
| `DEBUG` | Enable debug mode.                            | `false` |

## Running it locally

### 1. Start the go application

```bash
go run main.go
```
