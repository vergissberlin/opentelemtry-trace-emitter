const {NodeTracerProvider} = require('@opentelemetry/sdk-trace-node')
const {OTLPTraceExporter} = require('@opentelemetry/exporter-trace-otlp-grpc')
const {SimpleSpanProcessor} = require('@opentelemetry/sdk-trace-base')
const {trace, SpanStatusCode, context} = require('@opentelemetry/api')

const collectorUrl = process.env.OTEL_COLLECTOR_URL || 'http://localhost:4317'
const traceInterval = parseInt(process.env.TRACE_INTERVAL, 10) || 5000

const provider = new NodeTracerProvider()
const exporter = new OTLPTraceExporter({url: collectorUrl})

provider.addSpanProcessor(new SimpleSpanProcessor(exporter))
provider.register()

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

function generateRandomTrace() {
    const randomTraceName = traceNames[Math.floor(Math.random() * traceNames.length)]
    const randomProcessName = processNames[Math.floor(Math.random() * processNames.length)]

    const rootSpan = tracer.startSpan(randomTraceName)
    context.with(trace.setSpan(context.active(), rootSpan), () => {
        const span = tracer.startSpan(randomProcessName)
        span.setAttribute('random-value', Math.random())
        console.log(`Trace: ${randomTraceName} | Span: ${randomProcessName} | Trace ID: ${rootSpan.spanContext().traceId}`)
        span.setStatus({code: SpanStatusCode.OK})
        span.end()
    })

    rootSpan.end()
}

setInterval(generateRandomTrace, traceInterval)
