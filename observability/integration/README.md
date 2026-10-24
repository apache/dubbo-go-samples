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

Grafana is published on the standard host port `3000`. For a different host port, use a local Compose override.

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
dashboard automatically. Both applications enable the SDK's context logger with
`logger.WithTraceIntegration(true)`. Calls to `CtxInfof`/`CtxErrorf` include
`trace_id` and `span_id` only when the context contains a valid span.

## Capture and check the fields

Use an existing Nacos registry, then run:

```bash
python3 observability/integration/capture.py /tmp/dubbo-observability-capture
python3 observability/integration/validate.py /tmp/dubbo-observability-capture
python3 -m unittest discover -s observability/integration -p 'test_validate.py'
```

The output directory must be new. To use another registry, set
`DUBBO_OBSERVABILITY_REGISTRY_ADDRESS=127.0.0.1:18848`. The capture starts its own
Collector on ports `24318`–`24321` and stops only that container and its two
applications. Use `--port-base 25318` if these ports are occupied. It does not
start or stop Nacos, Jaeger, Prometheus or Grafana. Docker, Python 3 and Go are
required. Do not run another provider with these application and service names
in the same registry during the capture.

The script builds the existing provider and consumer, sends four successful
requests and one business error using group `observability` and version `1.0.0`,
and keeps the consumer alive long enough to scrape its metrics. It saves:

- `manifest.json`: SDK version, Go version, settings and exporter endpoints;
- `provider.prom` / `consumer.prom`: raw Prometheus responses;
- `traces.jsonl`: actual OTLP spans from the Collector file exporter;
- `provider.log` / `consumer.log`: application logs, including context logger calls;
- `collector.yaml` / `collector.log`: exporter configuration and diagnostics.

[capture/](capture/) checks in `manifest.json`, the Prometheus scrapes,
`traces.jsonl`, and `collector.yaml`. `provider.log`, `consumer.log`, and
`collector.log` stay in the `capture.py` output directory. The repository
ignores `*.log`, so those files are not part of the sample commit.
`validate.py` still requires the two application logs in the directory it
checks. The offline tests rebuild the context-log lines from `traces.jsonl`.

**Known result:** the pinned SDK currently labels consumer metrics with the
provider application name. The capture and validator exit 1 and report this
mismatch; the other checks still run. See [capture/](capture/) for the actual
output. This sample does not hide the mismatch with relabeling.

The validator checks both sides' counts (5 total, 4 successful, 1 failed), QPS
and RT series, non-empty group/version labels, and five CLIENT/SERVER span pairs.
For every pair it compares `service.name` with `application_name`, `rpc.service`
with `interface`, and `rpc.method` with `method`. The SERVER span must reference
the CLIENT span in the same trace. Both error spans must have ERROR status.
Context log IDs must exist in the exported spans; calls without a span or with
an invalid span context must omit those fields.

The SDK is pinned to merged commit `e4aabaf27773`. The earlier
`v3.3.2-20260910` release exports the composite service key as `rpc.service`, so
a non-empty group/version does not match the `interface` metric label. No SDK
changes are part of this sample. See [capture/](capture/) for a checked-in run
and the limits of the checks.

For Admin, configure Prometheus to scrape the two endpoints shown in the
manifest. Keep the capture applications running with the normal commands above
when inspecting live views. In **Application**, select each application and
compare total/success/failure counts and RT. In its service view, select
`observability.GreetService` / `Greet`, group `observability`, version `1.0.0`.
In the trace view, open one captured trace ID and check the two services and
parent relationship. Admin must use the same Prometheus and Jaeger backends as
this sample. The capture's file-only Collector is for automated checks; use the
normal Compose Collector/Jaeger stack for the Admin trace view.

These checks verify the exported data, not Admin rendering. Admin query changes
in [apache/dubbo-admin#1591](https://github.com/apache/dubbo-admin/pull/1591) and
metric refresh work in [#3764](https://github.com/apache/dubbo-go/pull/3764) /
[#3765](https://github.com/apache/dubbo-go/pull/3765) remain separate. A positive
QPS snapshot does not prove that the sliding window refreshes after traffic
stops. The sample does not require those changes to run.
