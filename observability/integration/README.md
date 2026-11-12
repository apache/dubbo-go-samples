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

### RPC latency histograms

This branch temporarily pins the SDK to the public commit from
[apache/dubbo-go#3791](https://github.com/apache/dubbo-go/pull/3791).
Replace that pin with an upstream SDK version containing the feature before merging
this sample. Older SDK versions do not provide the histogram options below.

Both applications enable cumulative latency histograms independently of
sliding-window aggregation:

```go
metrics.WithHistogramEnabled(),
metrics.WithHistogramBuckets([]float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.3, 1}),
```

Bucket bounds and observations are in **seconds**. These sample-specific bounds
include the provider's intentional 20 ms delay; omitting `WithHistogramBuckets`
uses the SDK defaults (0.1, 0.3, 0.5, 1, 3, 5, 10 seconds).
The exporter adds a `+Inf` bucket. Histograms record completed successful and
failed RPCs, with the existing application, IP, interface, method, group and
version labels; they do not add a `side` label.

Provider and consumer export these metric families respectively:

```text
dubbo_provider_rt_milliseconds_histogram_seconds_{bucket,sum,count}
dubbo_consumer_rt_milliseconds_histogram_seconds_{bucket,sum,count}
```

Despite the legacy `milliseconds` component of the names, the histogram unit is
seconds. The dashboard displays both latency distributions and aggregated p95/p99
estimates, rather than averaging per-instance percentiles. For example:

```promql
histogram_quantile(0.95,
  sum by (le) (rate(dubbo_provider_rt_milliseconds_histogram_seconds_bucket{job="dubbo-observability-server"}[1m]))
)
```

Run the continuous client above and allow at least two Prometheus scrapes before
inspecting rate-based panels. Quantiles are bucket-interpolated estimates; no
traffic or missing histograms must not be interpreted as zero latency.

### Validate histogram exports

With the sample's Nacos available, capture a finite run in a new directory:

```bash
python3 observability/integration/capture.py /tmp/dubbo-histogram-capture
python3 observability/integration/validate.py /tmp/dubbo-histogram-capture
python3 -m unittest discover -s observability/integration -p 'test_*.py'
```

New captures declare `histogram_buckets_seconds` in their manifest. The validator
checks provider and consumer counts against all five RPCs (four successes and one
failure), bucket bounds including `+Inf`, cumulative counts, matching RPC labels,
and a latency sum consistent with the buckets in seconds. The checked-in capture
predates histograms and remains a regression fixture for the original signals;
it is not evidence that this new histogram path ran successfully.

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
