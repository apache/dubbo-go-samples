# Dubbo Admin Integration

English | [中文](README_CN.md)

This sample runs Dubbo-go services under [Dubbo Admin](https://github.com/apache/dubbo-admin)
and walks through what Admin shows and controls for them: applications,
instances, services and method metadata, debug calls, condition routes, tag
routes, instance traffic disabling, metrics and traces.

```text
                       Dubbo Admin :38080
                    /        |          \
          registry/rules   PromQL      dashboards
                  |          |              |
               Nacos     Prometheus <---- Grafana :3000 ----> Jaeger
                  |          ^                                  ^
   consumer --> provider-a   |  scrape /metrics                 |
            \-> provider-b   +---- all three apps               |
                             OTLP traces --> OTel Collector ----+
```

## Layout

```text
admin/
├── go-server/cmd/main.go     # Provider; replies with its instance name
├── go-client/cmd/main.go     # Consumer; logs which instance served each call
├── verify/cmd/               # Checks everything below through Admin's API
├── proto/                    # adminsample.GreetService, version 1.0.0, group admin-sample
├── conf/                     # Dubbo Admin, Prometheus and Collector configuration
├── grafana/provisioning/     # Datasources with the UIDs Admin's dashboards expect
├── docker-compose.yaml
├── extract_dashboards.py     # Pulls Admin's Grafana dashboards out of its manifests
└── start.sh
```

The two providers run in separate containers because Admin disables traffic
per IP address.

## Run

Requirements: Go, Docker with Compose 2.17 or newer, Python 3 and curl.

```shell
cd admin
./start.sh
```

`start.sh` builds the sample binaries for Linux, downloads Admin's Grafana
dashboards, and starts the stack. No published Admin image contains the
features shown here yet, so the first start builds Admin from the commit pinned
in `.env`, which takes several minutes.

Then open:

| Service | Address |
| --- | --- |
| Dubbo Admin | http://localhost:38080/admin/ (user `admin`, password `dubbo@2025`) |
| Grafana | http://localhost:3000 |
| Prometheus | http://localhost:9090 |
| Jaeger | http://localhost:16686 |
| Nacos | http://localhost:8848/nacos |

Follow the consumer while you try the steps below:

```shell
docker compose logs -f consumer
```

Every second it sends one untagged call and one call tagged `gray`:

```text
untagged -> provider-a
tag=gray -> provider-b
```

To check every step at once, run:

```shell
docker compose run --rm verify
```

```text
PASS  Admin lists the provider application with two instances
PASS  Admin lists both provider instances with their IPs
PASS  Admin shows the service with its version and group
PASS  Admin shows the Greet method from published metadata
PASS  Admin invokes Greet on provider-a through its debug API
PASS  Without rules, calls reach both providers
PASS  Condition route from Admin pins calls to provider-a
PASS  Tag route from Admin sends gray calls to provider-b
PASS  Traffic disable from Admin drains provider-b
PASS  Admin queries consumer and provider request metrics
PASS  Traces reach Jaeger and Admin links the trace dashboard
all 11 checks passed
```

`verify` resets the rules it uses, so it replaces rules you created with the
same names and turns traffic on `provider-b` back on. Stop everything with
`docker compose down`. All ports are published on `127.0.0.1` only.

## What to look at in Admin

### Applications, instances and services

- **Resources > Applications** lists `admin-sample-provider` with two instances.
  With application-level service discovery, consumers do not register
  instances, so `admin-sample-consumer` does not appear in this list.
- The instance list shows the IPs of `provider-a` and `provider-b`. The routing
  steps below use them.
- **Resources > Services** lists `adminsample.GreetService` with version
  `1.0.0` and group `admin-sample`. Its method `Greet` takes
  `adminsample.GreetRequest` and returns `adminsample.GreetResponse`; the
  provider publishes these definitions to the Nacos metadata center.
- On the service's **Debug** tab, pick an instance and send
  `[{"name": "admin"}]`. Admin reads the message fields through the
  provider's Triple reflection service, which `dubbo-go/v3/imports` enables,
  and shows the reply with the instance that served it.

### Condition route

In **Traffic Management > Condition Rule**, create a rule for service
`adminsample.GreetService`, version `1.0.0`, group `admin-sample`, using the IP
of `provider-a`:

```yaml
configVersion: v3.0
scope: service
key: adminsample.GreetService
enabled: true
force: true
runtime: true
conditions:
  - "=> host = <provider-a IP>"
```

All calls in the consumer log now go to `provider-a`. Delete the rule to
restore balancing.

### Tag route

In **Traffic Management > Tag Rule**, create a rule for application
`admin-sample-provider`:

```yaml
configVersion: v3.0
scope: application
key: admin-sample-provider
enabled: true
force: false
runtime: true
tags:
  - name: gray
    addresses:
      - "<provider-b IP>:20000"
```

Calls tagged `gray` now go to `provider-b`, and untagged calls go to
`provider-a`, because untagged calls avoid addresses claimed by a tag. Under
application-level discovery, match tags by `addresses`; parameters set on the
provider are not visible to tag `match` conditions.

### Traffic disable

Open the `provider-b` instance and turn on **Flow disabled** in its
**Configuration** tab. The consumer log shows only
`provider-a`. Admin implements this as a condition rule
`admin-sample-provider.condition-router` with `=> host != <provider-b IP>`.
Turn the switch off to restore balancing.

### Metrics and traces

The **Monitor** and **Tracing** tabs of applications, instances and services
embed Admin's Grafana dashboards. Spans open in Grafana with links to the
request rate of the same RPC, filtered by the span's application: on a consumer
span the consumer link has data, on a provider span the provider link does.

Application-level panels also count the streaming calls each provider serves:
the Triple health-check watch a dubbo-go consumer opens to each provider when
it starts, which ends after about 20 seconds, and the reflection calls behind
Admin's debug page. They make the application's response time look high. The
service's **Monitor** tab shows `Greet` alone.

## Notes

- The consumer must configure a config center to receive Admin's rules, with
  group `dubbo`, which is where Admin writes them. Without it, the log reports
  `Config center does not start, ... router will not be enabled`.
- At startup the consumer logs
  `failed to query rule, key= group=dubbo ... dataId can not be empty`,
  because no data ID is configured for application config in the config
  center. Routing is not affected.
- While building this sample, a consumer that subscribed before the providers
  published the service name mapping sometimes found no providers for five
  minutes, and a rule added and removed within seconds of a client subscribing
  sometimes stayed applied in that client. The consumer therefore starts only
  after both providers serve, and `verify` waits a few seconds after
  subscribing before it changes rules.
- Admin's `Dockerfile` pulls its base images through
  `swr.cn-north-4.myhuaweicloud.com`. If that registry is unreachable, build the
  image yourself from the same commit with the Docker Hub images and tag it
  `dubbo-admin-sample:<DUBBO_ADMIN_REF>`; Compose then uses it.
- The Admin password and session secret in `conf/dubbo-admin.yaml` are for
  local use only.
- This sample is not part of the CI integration tests, which would have to
  build Admin on every run.
