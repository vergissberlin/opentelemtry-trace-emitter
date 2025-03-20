# opentelemtry-trace-emitter

If you want to test your OpenTelemtry configuration with a random span emitter.

## Usage

```bash
docker run -it --rm \
  -e OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4317 \
  -e OTEL_RESOURCE_ATTRIBUTES=service.name=example-service \
  ghcr.io/rafaelmartins/opentelemetry-trace-emitter:latest
```

## Environment Variables

* `OTEL_EXPORTER_OTLP_ENDPOINT`: The endpoint to send the spans to. Default is `http://localhost:4317`.
* `OTEL_RESOURCE_ATTRIBUTES`: The resource attributes to add to the spans. Default is `service.name=example-service`.
* `OTEL_TRACES_EXPORTER`: The exporter to use. Default is `otlp`.
* `OTEL_TRACES_EXPORTER_OTLP_ENDPOINT`: The endpoint to send the spans to. Default is `http://localhost:4317
* `OTEL_TRACES_EXPORTER_OTLP_INSECURE`: Whether to use insecure connection. Default is `false`.
* `OTEL_TRACES_EXPORTER_OTLP_HEADERS`: The headers to send with the spans. Default is `""`.
