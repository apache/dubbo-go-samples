# Dubbo Go Observability Integration Validation Sample

English | [中文](README_CN.md)

This is an integration-validation sample. It combines Dubbo Go's existing
observability capabilities into one runnable flow:

```text
Nacos registry
       |
Dubbo Triple client -> Dubbo Triple server
        |                    |
        +-- Trace --> OpenTelemetry Collector --> Jaeger
        +-- Metrics ---------------------------> Prometheus -> Grafana
        +-- Logs ------------------------------> trace_id / span_id in stdout
```

## Layout

```text
observability/
├── integration/                    # This end-to-end integration-validation sample
│   ├── go-client/cmd/main.go       # Triple consumer and request generator
│   ├── go-server/cmd/main.go       # Triple provider with a forced error path
│   ├── proto/                      # Self-contained service contract and stubs
│   ├── docker-compose.yaml         # Nacos, Collector, Jaeger, Prometheus, Grafana
│   ├── otel-collector-config.yaml  # OTLP Collector pipeline to Jaeger
│   ├── prometheus.yml              # Scrapes the two local metrics endpoints
│   └── grafana/                    # Provisioned datasource and dashboard
├── probe/                          # Kubernetes liveness/readiness/startup sample
└── prometheus_grafana/             # Push/Pull Metrics and Grafana sample
```

## Prerequisites

- Go 1.25 or newer;
- Docker with Compose;
- a local checkout of this repository.

The Compose stack is pinned to Nacos `v2.5.2`, OpenTelemetry Collector
`0.104.0`, Jaeger `1.57`, Prometheus `v2.55.1`, and Grafana `11.2.0` so that
the validation can be reproduced with a known toolchain.

For the repository-standard integration path, start the root dependencies and
run this sample through the root harness. The harness uses the root Makefile to
build and stop `go-server/cmd`, runs `go-client/cmd`, and starts the additional
observability services defined by this sample:

```bash
docker compose -f docker-compose.yml up -d
./integrate_test.sh observability/integration
```

The sample is also included in `start_integrate_test.sh`.

## Run the sample

Start Nacos, the OpenTelemetry Collector, Jaeger, Prometheus, and Grafana:

```bash
cd observability/integration
docker compose up -d
cd ../..
```

Wait for Nacos to become healthy before starting the applications:

```bash
docker compose -f observability/integration/docker-compose.yaml ps
```

The provider registers itself in Nacos and the consumer discovers it from the
registry. The default registry address is `127.0.0.1:8848`; override it with
`DUBBO_OBSERVABILITY_REGISTRY_ADDRESS` when the registry is elsewhere.

Grafana is published on the standard host port `3000`. If port `3000` is
already occupied in your local environment, stop the conflicting process; the
sample itself keeps the repository-wide standard port mapping.

In a second terminal, start the Dubbo provider:

```bash
go run ./observability/integration/go-server/cmd
```

In a third terminal, run the client. Every fifth request intentionally uses the
name `error` so that the failure path can be inspected in Metrics, Trace, and
logs:

```bash
go run ./observability/integration/go-client/cmd -requests 20 -interval 500ms
```

To keep generating traffic for dashboard exploration, pass `-requests 0`:

```bash
go run ./observability/integration/go-client/cmd -requests 0 -interval 500ms
```

To reproduce a timeout, set a shorter per-request deadline than the provider's
intentional 20 ms processing delay:

```bash
go run ./observability/integration/go-client/cmd -requests 5 -timeout 1ms
```

To reproduce an explicit cancellation rather than a deadline expiry:

```bash
go run ./observability/integration/go-client/cmd -requests 5 -cancel-after 1ms
```

The default client run is finite so that it can be used by the repository
integration harness. Pass `-requests 0` when continuous traffic is needed for
dashboard exploration.

Stop the client with `Ctrl-C`, stop the server with `Ctrl-C`, and remove the
monitoring stack when finished:

```bash
docker compose -f observability/integration/docker-compose.yaml down
```

## Inspect the signals

- Jaeger: <http://localhost:16686>
- Nacos: <http://localhost:8848/nacos>
- Prometheus: <http://localhost:9090>
- Grafana: <http://localhost:3000> (`admin` / `admin`)
- Provider metrics: <http://localhost:9099/prometheus>
- Consumer metrics: <http://localhost:9098/prometheus>

Grafana provisions the Prometheus datasource and the `Dubbo Go Observability`
dashboard automatically. The application logs remain in the provider and
client terminals; each application log reads the active OpenTelemetry
`SpanContext` and includes the same `trace_id` and `span_id` that can be found
in Jaeger. This sample deliberately uses the published module API. The
Dubbo-Go `CtxLogger` integration is covered by the dedicated
`logger/trace-integration` sample and is not reimplemented here.
