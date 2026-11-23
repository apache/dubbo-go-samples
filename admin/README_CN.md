# Dubbo Admin 对接示例

[English](README.md) | 中文

本示例让 Dubbo-go 服务接入 [Dubbo Admin](https://github.com/apache/dubbo-admin)，
逐项演示 Admin 能看到和控制的内容：应用、实例、服务与方法元数据、服务调试、
条件路由、标签路由、实例禁流、指标和链路追踪。

```text
                       Dubbo Admin :38080
                    /        |          \
          注册/规则        PromQL        面板
                  |          |              |
               Nacos     Prometheus <---- Grafana :3000 ----> Jaeger
                  |          ^                                  ^
   consumer --> provider-a   |  抓取 /metrics                   |
            \-> provider-b   +---- 三个应用                     |
                             OTLP trace --> OTel Collector -----+
```

## 目录结构

```text
admin/
├── go-server/cmd/main.go     # Provider，响应中带上自己的实例名
├── go-client/cmd/main.go     # Consumer，日志显示每次调用由哪个实例处理
├── verify/cmd/               # 通过 Admin API 检查下文的全部内容
├── proto/                    # adminsample.GreetService，版本 1.0.0，分组 admin-sample
├── conf/                     # Dubbo Admin、Prometheus 和 Collector 配置
├── grafana/provisioning/     # 数据源 UID 与 Admin 面板引用的一致
├── docker-compose.yaml
├── extract_dashboards.py     # 从 Admin 的部署清单中提取 Grafana 面板
└── start.sh
```

两个 provider 分别运行在不同容器中，因为 Admin 按 IP 禁用流量。

## 运行

依赖：Go、Docker（Compose 2.17 或更高版本）、Python 3 和 curl。

```shell
cd admin
./start.sh
```

`start.sh` 会编译 Linux 版本的示例程序、下载 Admin 的 Grafana 面板并启动整套环境。
目前还没有包含这些功能的 Admin 发布镜像，所以首次启动时会按 `.env` 中固定的 commit
从源码构建 Admin，需要几分钟。

启动后可以访问：

| 服务 | 地址 |
| --- | --- |
| Dubbo Admin | http://localhost:38080/admin/（用户名 `admin`，密码 `dubbo@2025`） |
| Grafana | http://localhost:3000 |
| Prometheus | http://localhost:9090 |
| Jaeger | http://localhost:16686 |
| Nacos | http://localhost:8848/nacos |

操作下面的步骤时，可以同时观察 consumer 日志：

```shell
docker compose logs -f consumer
```

consumer 每秒发出一次不带标签的调用和一次带 `gray` 标签的调用：

```text
untagged -> provider-a
tag=gray -> provider-b
```

一次性检查全部步骤：

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

`verify` 会重置它用到的规则，因此会覆盖你手动创建的同名规则，并重新打开 `provider-b` 的流量。
用 `docker compose down` 停止全部服务。所有端口只发布在 `127.0.0.1` 上。

## 在 Admin 中查看

### 应用、实例和服务

- **资源详情 > 应用** 中列出 `admin-sample-provider`，共两个实例。使用应用级服务发现时，
  consumer 不注册实例，所以列表中没有 `admin-sample-consumer`。
- 实例列表中可以看到 `provider-a` 和 `provider-b` 的 IP，下面的路由步骤会用到。
- **资源详情 > 服务** 中列出 `adminsample.GreetService`，版本 `1.0.0`，分组
  `admin-sample`。方法 `Greet` 的参数为 `adminsample.GreetRequest`，返回
  `adminsample.GreetResponse`，这些定义由 provider 发布到 Nacos 元数据中心。
- 在服务的 **调试** 标签页选择一个实例，发送 `[{"name": "admin"}]`。Admin 通过 provider 的
  Triple 反射服务（导入 `dubbo-go/v3/imports` 时默认开启）获取消息字段，并显示响应以及处理它的实例。

### 条件路由

在 **流量管控 > 条件路由** 中，为服务 `adminsample.GreetService`（版本 `1.0.0`，
分组 `admin-sample`）创建规则，填入 `provider-a` 的 IP：

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

consumer 日志中的调用全部转到 `provider-a`。删除规则后恢复负载均衡。

### 标签路由

在 **流量管控 > 标签路由** 中，为应用 `admin-sample-provider` 创建规则：

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

带 `gray` 标签的调用转到 `provider-b`，不带标签的调用转到 `provider-a`，
因为不带标签的调用会避开已被标签占用的地址。使用应用级服务发现时，请用 `addresses`
匹配实例；provider 上设置的参数对标签的 `match` 条件不可见。

### 实例禁流

打开 `provider-b` 实例，在 **场景配置** 标签页中打开 **流量禁用**。consumer 日志中只剩
`provider-a`。Admin 的实现方式是写入条件路由 `admin-sample-provider.condition-router`，
内容为 `=> host != <provider-b IP>`。关闭开关后恢复负载均衡。

### 指标和链路追踪

应用、实例和服务的 **监控** 和 **链路追踪** 标签页内嵌了 Admin 的 Grafana 面板。
在 Grafana 中打开 span，可以跳转到同一 RPC 的请求速率，并按该 span 所属应用过滤：
consumer 的 span 上 consumer 链接有数据，provider 的 span 上 provider 链接有数据。

应用维度的面板还会统计 provider 处理的流式调用：dubbo-go consumer 启动时会向每个 provider
发起一条 Triple 健康检查流（约 20 秒后结束），Admin 调试页也会发起反射调用。它们会让应用的响应时间显得偏高。
只看 `Greet` 请查看服务的 **监控** 标签页。

## 说明

- consumer 必须配置配置中心才能收到 Admin 下发的规则，并且分组设为 `dubbo`，
  即 Admin 写入规则的分组。否则日志会提示
  `Config center does not start, ... router will not be enabled`。
- consumer 启动时会输出
  `failed to query rule, key= group=dubbo ... dataId can not be empty`。
  这是因为没有为配置中心设置应用配置的 data ID，不影响路由。
- 编写本示例时观察到：provider 发布服务名映射之前就订阅的 consumer，有时会在 5 分钟内找不到 provider；
  客户端订阅后几秒内新增又删除的规则，有时会在该客户端中一直生效。因此 consumer 会等两个 provider
  都开始提供服务后才启动，`verify` 订阅后也会先等几秒再修改规则。
- Admin 的 `Dockerfile` 通过 `swr.cn-north-4.myhuaweicloud.com` 拉取基础镜像。
  如果该镜像源不可达，可以用同一 commit 和 Docker Hub 上的镜像自行构建，
  并打上 `dubbo-admin-sample:<DUBBO_ADMIN_REF>` 标签，Compose 会直接使用它。
- `conf/dubbo-admin.yaml` 中的 Admin 密码和 session secret 仅供本地使用。
- 本示例没有加入 CI 集成测试，否则每次运行都要构建 Admin。
