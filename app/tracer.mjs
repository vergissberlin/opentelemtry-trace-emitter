import * as opentelemetry from '@opentelemetry/sdk-node'
import {getNodeAutoInstrumentations} from '@opentelemetry/auto-instrumentations-node'
import {OTLPMetricExporter} from '@opentelemetry/exporter-metrics-otlp-proto'
import {PeriodicExportingMetricReader} from '@opentelemetry/sdk-metrics'
import {CollectorTraceExporter} from '@opentelemetry/exporter-collector-grpc'
import {Resource} from '@opentelemetry/resources'
import {
    SEMRESATTRS_SERVICE_NAME,
    SEMRESATTRS_SERVICE_VERSION,
    SemanticAttributes
} from '@opentelemetry/semantic-conventions'
import {trace, context, SpanStatusCode} from '@opentelemetry/api'
import {MeterProvider} from '@opentelemetry/sdk-metrics'

const collectorUrl = process.env.OTEL_COLLECTOR_ENDPOINT || 'http://localhost:4317'
const traceInterval = parseInt(process.env.TRACE_INTERVAL, 10) || 5000

const sdk = new opentelemetry.NodeSDK({
    resource: new Resource({
        [SEMRESATTRS_SERVICE_NAME]: 'opentelemetry-trace-emitter',
        [SEMRESATTRS_SERVICE_VERSION]: '1.0',
        [SemanticAttributes.SERVICE_INSTANCE_ID]: '1',
    }),
    traceExporter: new CollectorTraceExporter({
        url: collectorUrl
    }),
    metricReader: new PeriodicExportingMetricReader({
        exporter: new OTLPMetricExporter({
            url: `${collectorUrl}/v1/metrics`,
            headers: {},
        }),
    }),
    instrumentations: [getNodeAutoInstrumentations()],
})
sdk.start()

const meterProvider = new MeterProvider({
    resource: new Resource({
        [SEMRESATTRS_SERVICE_NAME]: 'opentelemetry-trace-emitter',
        [SEMRESATTRS_SERVICE_VERSION]: '1.0',
        [SemanticAttributes.SERVICE_INSTANCE_ID]: '1',
    }),
    exporter: new OTLPMetricExporter({
        url: `${collectorUrl}/v1/metrics`,
        headers: {},
    }),
    interval: 1000,
})

const meter = meterProvider.getMeter('custom-metrics')
const requestCount = meter.createCounter('request_count', {
    description: 'Count of requests',
})
const requestDuration = meter.createHistogram('request_duration', {
    description: 'Duration of requests',
})

const tracer = trace.getTracer('random-tracer')

const processNames = [
    'database-query',
    'user-authentication',
    'http-request',
    'cache-fetch',
    'file-read',
    'message-processing',
    'background-task',
    'api-call',
    'data-validation',
    'report-generation'
]

const traceNames = [
    'order-processing',
    'user-login-flow',
    'checkout-session',
    'data-synchronization',
    'background-cleanup',
    'inventory-update',
    'payment-processing',
    'email-notification',
    'file-upload',
    'system-monitoring'
]

const randomAttributes = [
    {key: SemanticAttributes.ENDUSER_ID, value: () => `user-${Math.floor(Math.random() * 1000)}`},
    {key: SemanticAttributes.HTTP_STATUS_CODE, value: () => [200, 201, 400, 401, 403, 500][Math.floor(Math.random() * 6)]},
    {key: SemanticAttributes.DB_STATEMENT, value: () => ['SELECT * FROM users', 'UPDATE orders SET status="shipped"', 'DELETE FROM sessions'][Math.floor(Math.random() * 3)]},
    {key: 'cache.hit', value: () => Math.random() > 0.5},
    {key: 'processing.time_ms', value: () => Math.floor(Math.random() * 500)}
]

const randomLogs = [
    'Starting process',
    'Fetching data',
    'Processing request',
    'Sending response',
    'Operation completed',
    'Error occurred',
    'Retrying operation',
    'Operation successful',
    'Data saved',
    'Connection established'
]

function generateRandomTrace() {
    const randomTraceName = traceNames[Math.floor(Math.random() * traceNames.length)]
    const randomProcessName = processNames[Math.floor(Math.random() * processNames.length)]

    const rootSpan = tracer.startSpan(randomTraceName)
    context.with(trace.setSpan(context.active(), rootSpan), () => {
        const span = tracer.startSpan(randomProcessName)

        randomAttributes.forEach(attr => {
            if (Math.random() > 0.5) {
                span.setAttribute(attr.key, attr.value())
            }
        })

        const randomLog = randomLogs[Math.floor(Math.random() * randomLogs.length)]
        span.addEvent(randomLog)

        console.log(`Trace: ${randomTraceName} | Span: ${randomProcessName} | Trace ID: ${rootSpan.spanContext().traceId} | Log: ${randomLog}`)
        span.setStatus({code: SpanStatusCode.OK})
        span.end()

        // Record custom metrics
        requestCount.add(1, {process: randomProcessName})
        requestDuration.record(Math.random() * 1000, {process: randomProcessName})
    })

    rootSpan.end()
}

setInterval(generateRandomTrace, traceInterval)
