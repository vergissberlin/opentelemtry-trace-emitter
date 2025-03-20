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

const collectorUrl = process.env.OTEL_COLLECTOR_URL || 'http://localhost:4317'
const traceInterval = parseInt(process.env.TRACE_INTERVAL, 10) || 5000

const sdk = new opentelemetry.NodeSDK({
    resource: new Resource({
        [SEMRESATTRS_SERVICE_NAME]: 'opentelemetry-trace-emitter',
        [SEMRESATTRS_SERVICE_VERSION]: '1.0',
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

        console.log(`Trace: ${randomTraceName} | Span: ${randomProcessName} | Trace ID: ${rootSpan.spanContext().traceId}`)
        span.setStatus({code: SpanStatusCode.OK})
        span.end()
    })

    rootSpan.end()
}

setInterval(generateRandomTrace, traceInterval)
