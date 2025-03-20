FROM node:20-alpine

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm install

COPY . .

ENV OTEL_COLLECTOR_URL=http://localhost:4318/v1/traces
ENV TRACE_INTERVAL=5000

CMD ["node", "tracer.js"]
