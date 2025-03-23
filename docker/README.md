# Test with Docker Compose

## 1. Start the Docker container

```bash
docker network create opentelemetry-emitter-network
docker compose -f docker/compose.yaml down
docker compose -f docker/compose.yaml pull
docker compose -f docker/compose.yaml up -d
```

## 2. Access the visualisation application

Take a look to the logs of the visualisation application to get the access URL.
You can ope it by clicking on the link in the logs while holding the CTRL key.

```bash
docker compose -f docker/compose.yaml logs aspire 
```

## 3. Stop the Docker container

```bash
docker compose -f docker/compose.yaml ps
```

```bash
docker compose -f docker/compose.yaml logs -f collector
```

## 3. Local build

```bash
docker build -t opentelemetry-trace-emitter ./app/
```
