# Dubbo Go 可观测性集成验证样例

[English](README.md) | 中文

这个样例是一个集成验证样例，把 Dubbo Go 现有的可观测性能力组合成一条可运行链路：

```text
Nacos 注册中心
       |
Dubbo Triple 客户端 -> Dubbo Triple 服务端
        |                    |
        +-- Trace --> OpenTelemetry Collector --> Jaeger
        +-- Metrics ---------------------------> Prometheus -> Grafana
        +-- 日志 ------------------------------> stdout 中的 trace_id / span_id
```

## 目录结构

```text
observability/
├── integration/                    # 本端到端集成验证样例
│   ├── go-client/cmd/main.go       # Triple 消费端和请求生成器
│   ├── go-server/cmd/main.go       # Triple 提供端，包含可复现错误路径
│   ├── proto/                      # 独立的服务契约和生成代码
│   ├── docker-compose.yaml         # Nacos、Collector、Jaeger、Prometheus、Grafana
│   ├── otel-collector-config.yaml  # 转发到 Jaeger 的 OTLP Collector 管道
│   ├── prometheus.yml              # 抓取两个本地 Metrics endpoint
│   └── grafana/                    # 自动配置数据源和 Dashboard
├── probe/                          # Kubernetes 存活、就绪和启动探针样例
└── prometheus_grafana/             # Push/Pull Metrics 和 Grafana 样例
```

## 前置条件

- Go 1.25 或更新版本；
- Docker 和 Compose；
- 本仓库的本地 checkout。

Compose 中固定了 Nacos `v2.5.2`、OpenTelemetry Collector `0.104.0`、Jaeger
`1.57`、Prometheus `v2.55.1` 和 Grafana `11.2.0`，以便使用确定的工具链复现
这条验证链路。

按仓库标准进行集成验证时，先启动根目录依赖，再使用根级集成测试入口。该入口会
通过根目录 Makefile 构建和停止 `go-server/cmd`，运行 `go-client/cmd`，并启动本样例
额外声明的可观测性服务：

```bash
docker compose -f docker-compose.yml up -d
./integrate_test.sh observability/integration
```

本样例也已经加入 `start_integrate_test.sh` 的集成测试列表。

## 启动样例

启动 Nacos、OpenTelemetry Collector、Jaeger、Prometheus 和 Grafana：

```bash
cd observability/integration
docker compose up -d
cd ../..
```

启动应用前，先等待 Nacos 健康：

```bash
docker compose -f observability/integration/docker-compose.yaml ps
```

服务端会注册到 Nacos，客户端从 Nacos 发现服务。默认注册中心地址是
`127.0.0.1:8848`；如果注册中心在其他主机，可以通过
`DUBBO_OBSERVABILITY_REGISTRY_ADDRESS` 覆盖。

Grafana 默认映射到标准宿主机 `3000` 端口。如果你的本地环境已经占用
`3000`，请停止冲突进程；样例本身保持仓库统一的标准端口映射。

在第二个终端启动 Dubbo 服务端：

```bash
go run ./observability/integration/go-server/cmd
```

在第三个终端启动客户端。客户端每 5 个请求会故意发送一个 `error` 请求，
用于在 Metrics、Trace 和日志中检查失败路径：

```bash
go run ./observability/integration/go-client/cmd -requests 20 -interval 500ms
```

如果希望持续产生流量以观察 Dashboard，请传入 `-requests 0`：

```bash
go run ./observability/integration/go-client/cmd -requests 0 -interval 500ms
```

如果要复现超时，可以设置一个小于服务端故意延迟 20 毫秒的请求超时：

```bash
go run ./observability/integration/go-client/cmd -requests 5 -timeout 1ms
```

如果要复现显式取消而不是 deadline 到期：

```bash
go run ./observability/integration/go-client/cmd -requests 5 -cancel-after 1ms
```

客户端默认执行有限次数，这样可以直接接入仓库集成测试；如果希望持续产生流量观察
Dashboard，请显式传入 `-requests 0`。

测试完成后使用 `Ctrl-C` 停止客户端和服务端，再回收监控栈：

```bash
docker compose -f observability/integration/docker-compose.yaml down
```

### RPC 延迟 Histogram

客户端和服务端均开启累计延迟 Histogram，不依赖滑动窗口聚合：

```go
metrics.WithHistogramEnabled(),
metrics.WithHistogramBuckets([]float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.3, 1}),
```

桶边界和记录值的单位都是**秒**。样例选择的桶能区分服务端故意加入的 20 ms 延迟；
省略 `WithHistogramBuckets` 时，SDK 使用默认桶 0.1、0.3、0.5、1、3、5、10 秒。
Exporter 自动添加 `+Inf` 桶。成功和失败的已完成 RPC 都会记录，标签沿用应用名、IP、
接口、方法、group 和 version，不额外添加 `side` 标签。

服务端和客户端分别导出：

```text
dubbo_provider_rt_milliseconds_histogram_seconds_{bucket,sum,count}
dubbo_consumer_rt_milliseconds_histogram_seconds_{bucket,sum,count}
```

名称中的 `milliseconds` 是历史命名，不能据此把 Histogram 当成毫秒。
Dashboard 展示两侧的延迟分布及聚合 P95/P99 估计值，不再平均各实例的分位数。例如：

```promql
histogram_quantile(0.95,
  sum by (le) (rate(dubbo_provider_rt_milliseconds_histogram_seconds_bucket{job="dubbo-observability-server"}[1m]))
)
```

使用前面的持续调用命令，等待 Prometheus 至少抓取两次，再查看基于 rate 的面板。
分位数通过桶内插值估算；停流或未采到 Histogram 不代表延迟为零。

### 验证 Histogram 导出

客户端在前 5 次调用（4 次成功、1 次业务失败）后暂停发请求，等待两端异步指标更新，
检查计数增量、包含 `+Inf` 的配置桶、累计计数、RPC 标签和秒单位的延迟总和。
校验失败会 `panic`，使集成测试失败。少于 5 次的有限调用按实际次数校验；
超时、主动取消实验不执行这项校验。

标准集成脚本会等到客户端校验成功后，才停止持续调用。服务端指标地址不是默认值时，
为客户端传入 `-provider-metrics-url http://host:port/prometheus`。

原有捕获工具仍可用于采集有限调用的 Trace、Metric 与 Log 关联数据：

```bash
python3 observability/integration/capture.py /tmp/dubbo-histogram-capture
python3 observability/integration/validate.py /tmp/dubbo-histogram-capture
```

Histogram 校验由 Go 客户端执行，不依赖单独运行 Python 单测。仓库已有的快照采集于
Histogram 功能之前，继续用于原有信号的回归测试，不能作为 Histogram 已实跑通过的证明。

## 查看信号

- Jaeger：<http://localhost:16686>
- Nacos：<http://localhost:8848/nacos>
- Prometheus：<http://localhost:9090>
- Grafana：<http://localhost:3000>（`admin` / `admin`）
- 服务端 Metrics：<http://localhost:9099/prometheus>
- 客户端 Metrics：<http://localhost:9098/prometheus>

Grafana 会自动配置 Prometheus 数据源和 `Dubbo Go Observability` Dashboard。
两个应用均通过 `logger.WithTraceIntegration(true)` 启用 SDK 的上下文日志。
`CtxInfof` / `CtxErrorf` 在 context 含有效 span 时自动写入 `trace_id`、`span_id`，
没有有效 span 时不输出这些字段。
