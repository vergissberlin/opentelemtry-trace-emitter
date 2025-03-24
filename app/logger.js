import * as opentelemetry from '@opentelemetry/api'
import winston from 'winston'

const tracingFormat = function () {
    return winston.format(info => {
        const span = opentelemetry.trace.getSpan(opentelemetry.context.active())
        if (span) {
            const {spanId, traceId} = span.spanContext()
            const traceIdEnd = traceId.slice(traceId.length / 2)
            info['dd.trace_id'] = BigInt(`0x${traceIdEnd}`).toString()
            info['dd.span_id'] = BigInt(`0x${spanId}`).toString()
        }
        return info
    })()
}

const logger = winston.createLogger({
    transports: [new winston.transports.Console()],
    format: winston.format.combine(tracingFormat(), winston.format.json())
})

export default logger
