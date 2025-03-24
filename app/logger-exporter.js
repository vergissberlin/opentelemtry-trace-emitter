import logs from '@opentelemetry/api'
import {LoggerProvider, SimpleLogRecordProcessor} from '@opentelemetry/sdk-logs'
import OTLPLogExporter from '@opentelemetry/exporter-otlp-grpc'

// Log-Exporter erstellen
const logExporter = new OTLPLogExporter({
    url: 'http://localhost:4317/v1/logs' // oder 'http://localhost:4317' für gRPC
})

// Logger konfigurieren
const loggerProvider = new LoggerProvider()
loggerProvider.addLogRecordProcessor(new SimpleLogRecordProcessor(logExporter))
logs.setGlobalLoggerProvider(loggerProvider)

// Logger für spätere Nutzung abrufen
const logger = logs.getLogger('opentelemetry-trace-emitter')

// Beispiel-Log schreiben
logger.emit({
    severityNumber: logs.SeverityNumber.INFO,
    severityText: 'INFO',
    body: 'Dies ist eine OpenTelemetry-Log-Nachricht',
})

export {loggerProvider, logger}
