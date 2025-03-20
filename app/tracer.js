const {NodeTracerProvider} = require('@opentelemetry/sdk-trace-node')
const {OTLPTraceExporter} = require('@opentelemetry/exporter-trace-otlp-http')
const {SimpleSpanProcessor} = require('@opentelemetry/sdk-trace-base')
const {trace, SpanStatusCode} = require('@opentelemetry/api')

const collectorUrl = process.env.OTEL_COLLECTOR_URL || 'http://localhost:4318/v1/traces'
const traceInterval = parseInt(process.env.TRACE_INTERVAL, 10) || 5000

const provider = new NodeTracerProvider()
const exporter = new OTLPTraceExporter({url: collectorUrl})

provider.addSpanProcessor(new SimpleSpanProcessor(exporter))
provider.register()

const tracer = trace.getTracer('random-tracer')

function generateRandomTrace() {
    const span = tracer.startSpan('random-span')
    span.setAttribute('random-value', Math.random())
    console.log(`Trace created: ${span.spanContext().traceId}`)
    span.setStatus({code: SpanStatusCode.OK})
    span.end()
}

setInterval(generateRandomTrace, traceInterval)
