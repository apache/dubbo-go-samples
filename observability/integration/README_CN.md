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
两个应用均通过 `logger.WithTraceIntegration(true)` 启用 SDK 的上下文日志。
`CtxInfof` / `CtxErrorf` 在 context 含有效 span 时自动写入 `trace_id`、`span_id`，
没有有效 span 时不输出这些字段。

## 采集并检查字段

先准备可用的 Nacos，再从仓库根目录运行：

```bash
python3 observability/integration/capture.py /tmp/dubbo-observability-capture
python3 observability/integration/validate.py /tmp/dubbo-observability-capture
python3 -m unittest discover -s observability/integration -p 'test_validate.py'
```

输出目录必须不存在。注册中心地址可通过
`DUBBO_OBSERVABILITY_REGISTRY_ADDRESS=127.0.0.1:18848` 指定。
脚本只启动并清理本次采集的 Collector 和应用进程，不操作已有的注册中心或监控服务。
默认占用四个连续端口 `24318`–`24321`，可用 `--port-base 25318` 修改。
需要 Docker、Python 3 和 Go。同一注册中心中请勿同时运行同名应用和服务的其他 provider。

采集使用非空 group `observability` 和 version `1.0.0`，发出四次成功请求和一次业务错误。
目录中包含原始 Prometheus 响应、OTLP spans、应用日志、Collector 配置及 SDK/Go 版本和端点清单。
校验覆盖双方计数、QPS/RT、metric 与 trace 字段对应关系、五对 CLIENT/SERVER 的父子关系和错误状态，
以及 CtxLogger 在有效、无和无效 trace context 下的输出。真实采集节选见 [capture/](capture/)。

**当前结果：** SDK 导出的 consumer metric `application_name` 实际为 provider 应用名，
与 consumer 的 `service.name` 不同。采集脚本和校验脚本会返回退出码 1 并报告差异，
其余检查仍会执行。样例不通过重写标签或把两进程设成同名来掩盖它；这需要独立运行时修复。

本样例通过 `logger.WithTraceIntegration(true)` 启用 SDK 的上下文日志，使用 `CtxInfof` / `CtxErrorf`。
SDK 固定在已合并提交 `e4aabaf27773`；旧版 `v3.3.2-20260910` 的 `rpc.service` 含 group/version，
无法与 metric 的 `interface` 对齐。因此只更新 SDK 的固定版本，不在样例中修改 SDK 行为。

检查 Admin 页面时，使用原有 Compose（含 Jaeger），将 Admin 指向同一 Prometheus/Jaeger。
按 manifest 中的端点抓取指标，在应用页检查两侧计数和 RT；在服务页选择
`observability.GreetService` / `Greet` 及上述 group/version；在 trace 页用采集到的 trace ID
核对双方服务和父子关系。自动采集使用仅写文件的 Collector，不提供 Jaeger 查询接口。
这些自动检查只验证导出的数据，不代表已经验证 Admin 页面。
Admin #1591 的查询调整与 dubbo-go #3764/#3765 的滑动窗口刷新仍独立推进，
本样例不依赖它们；一次 QPS 快照不能证明停止流量后指标会刷新。
