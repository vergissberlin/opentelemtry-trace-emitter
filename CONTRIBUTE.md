# Contribute

## Run the script

```bash
cd app
npm start
```

## Build the image

```bash
docker build -t opentelemetry-trace-emitter .
```

## Test the image

```bash
docker run -it --rm \
    -e OTEL_COLLECTOR_ENDPOINT=http://localhost:4317 \
    opentelemetry-trace-emitter
```
