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

## 查看信号

- Jaeger：<http://localhost:16686>
- Nacos：<http://localhost:8848/nacos>
- Prometheus：<http://localhost:9090>
- Grafana：<http://localhost:3000>（`admin` / `admin`）
- 服务端 Metrics：<http://localhost:9099/prometheus>
- 客户端 Metrics：<http://localhost:9098/prometheus>

Grafana 会自动配置 Prometheus 数据源和 `Dubbo Go Observability` Dashboard。
应用日志保留在客户端和服务端终端中；每条应用日志都会读取当前 OpenTelemetry
`SpanContext`，并包含和 Jaeger 中相同的 `trace_id`、`span_id`。这个样例刻意只使用
已发布的 module API；Dubbo-Go 的 `CtxLogger` 集成由专门的
`logger/trace-integration` 样例覆盖，这里不重复实现。
