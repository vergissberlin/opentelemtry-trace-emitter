FROM golang AS builder


WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o app .
RUN chmod +x app

FROM alpine:latest AS runner

RUN apk update && \
    apk add --no-cache --repository=http://dl-cdn.alpinelinux.org/alpine/edge/testing grpcurl


WORKDIR /app
COPY --from=builder /app/app .

ENV OTEL_COLLECTOR_ENDPOINT=http://localhost:4318/v1/traces \
    TRACE_INTERVAL=5000 \
    SPAN_COUNT_MAX=12 \
    DEBUG=true

CMD ["./app"]
