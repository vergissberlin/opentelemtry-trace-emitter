FROM node:22-alpine

RUN apk update && \
    apk add --no-cache --repository=http://dl-cdn.alpinelinux.org/alpine/edge/testing grpcurl

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm install

COPY . .

ENV OTEL_COLLECTOR_ENDPOINT=http://localhost:4318/v1/traces \
    TRACE_INTERVAL=5000 \
    SPAN_COUNT_MAX=12 \
    DEBUG=true

CMD ["node", "index.js"]
