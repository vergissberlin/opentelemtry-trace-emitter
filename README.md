# opentelemtry-trace-emitter

If you want to test your OpenTelemtry configuration with a random span emitter.

## Usage

### 1. Pull the image

```bash
docker pull ghcr.io/vergissberlin/opentelemetry-trace-emitter:latest
```

### 2. Run the image

```bash
docker run -it --rm \
    -e OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4317 \
    ghcr.io/vergissberlin/opentelemetry-trace-emitter
```

## Environment Variables

|         Environment Variable         |                 Description                  |            Default             |
| ------------------------------------ | -------------------------------------------- | ------------------------------ |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | The endpoint to send the spans to.           | `http://localhost:4317` |
| `OTEL_RESOURCE_ATTRIBUTES` | The resource attributes to add to the spans. | `service.name=example-service` |
| `OTEL_TRACES_EXPORTER` | The exporter to use.                         | `otlp` |
| `OTEL_TRACES_EXPORTER_OTLP_ENDPOINT` | The endpoint to send the spans to.           | `http://localhost:4317` |
| `OTEL_TRACES_EXPORTER_OTLP_INSECURE` | Whether to use insecure connection.          | `false` |
| `OTEL_TRACES_EXPORTER_OTLP_HEADERS` | The headers to send with the spans.          | `""` |
