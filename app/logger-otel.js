import * as api from '@opentelemetry/api'
const {logs} = api
import {LoggerProvider, BatchLogRecordProcessor} from '@opentelemetry/sdk-logs'
import {OTLPLogExporter} from '@opentelemetry/exporter-logs-otlp-grpc'
import {Resource} from '@opentelemetry/resources'
import {SemanticResourceAttributes} from '@opentelemetry/semantic-conventions'

// OTLP HTTP Log Exporter konfigurieren
const logExporter = new OTLPLogExporter({
    url: 'http://otel-collector:4318/v1/logs', // Stelle sicher, dass dein OTel Collector diesen Port nutzt
})

// Logger Provider mit Exporter erstellen
const loggerProvider = new LoggerProvider({
    resource: new Resource({
        [SemanticResourceAttributes.SERVICE_NAME]: 'mein-node-service',
        [SemanticResourceAttributes.SERVICE_INSTANCE_ID]: 'instance-1',
    }),
})

loggerProvider.addLogRecordProcessor(new BatchLogRecordProcessor(logExporter))

// OpenTelemetry Logger abrufen
const logger = loggerProvider.getLogger('otel-logger')

// Logging-Funktion
class Logger {
    log(level, message) {
        const time = new Date().toISOString()

        logger.emit({
            severityText: level.toUpperCase(),
            body: message,
            attributes: {
                'service.name': 'mein-node-service',
                'service.instance.id': 'instance-1',
                'timestamp': time,
            },
        })

        console.log(JSON.stringify({time, level, message}))
    }
}

// Den Logger exportieren
export default new Logger()
