# opentelemtry-trace-emitter

If you want to test your OpenTelemtry configuration with a random span emitter.

## Usage

### 1. Pull the images

```bash
docker pull ghcr.io/vergissberlin/opentelemetry-trace-emitter:javascript
docker pull ghcr.io/vergissberlin/opentelemetry-trace-emitter:go
```

### 2. Run the images

```bash
docker run -it --rm \
    -e OTEL_COLLECTOR_ENDPOINT=http://localhost:4317 \
    ghcr.io/vergissberlin/opentelemetry-trace-emitter:javascript
```

```bash
docker run -it --rm \
    -e OTEL_COLLECTOR_ENDPOINT=http://localhost:4317 \
    ghcr.io/vergissberlin/opentelemetry-trace-emitter:go
```

## Environment Variables

|   Environment Variable    |                  Description                   |         Default         |
| ------------------------- | ---------------------------------------------- | ----------------------- |
| `OTEL_COLLECTOR_ENDPOINT` | The endpoint to send the spans to.             | `http://localhost:4317` |
| `TRACE_INTERVAL` | The timeout for the exporter.                  | `5000` |
| `SPAN_COUNT_MAX` | The maximum number of spans per trace to send. | `12` |
| `DEBUG` | Enable debug mode.                             | `false` |
