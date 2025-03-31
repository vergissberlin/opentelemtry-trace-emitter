# Contribute

## Run the script

```bash
cd app
npm start
```

## Build the image

```bash
docker build -t opentelemetry-trace-emitter:javascript ./app/javascript
docker build -t opentelemetry-trace-emitter:go ./app/go
```

## Test the javascript image

```bash
docker run -it --rm \
    -e OTEL_COLLECTOR_ENDPOINT=http://localhost:4317 \
    opentelemetry-trace-emitter:javascript
```

### Test go image

```bash
docker run -it --rm \
    -e OTEL_COLLECTOR_ENDPOINT=http://localhost:4317 \
    opentelemetry-trace-emitter:go
```
