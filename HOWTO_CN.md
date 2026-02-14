# 如何运行

[English](HOWTO.md) | 中文

目前有两种方式来运行 dubbo-go 的示例:

1. 脚本驱动的集成测试（CI 流程）
2. 不使用脚本的手动运行

## 1. 集成测试（CI 流程）

### 1.1 集成测试流程

当前 CI 使用脚本驱动的集成测试。

对单个 sample，`./integrate_test.sh <sample-path>` 的执行顺序如下：
1. 启动 `go-server`
2. 启动辅助 Go server（`*-server/cmd/*.go`，如 `grpc-server`）
3. 运行 `go-client`
4. 运行 `java-client`（如果存在）
5. 停止 `go-server`
6. 启动 `java-server`（如果存在）
7. 运行 `java-client`（如果存在）
8. 再次运行 `go-client`（验证 Go client 能调用 Java server）

说明：
- 如果环境中没有 `mvn`，会自动跳过 Java 阶段。
- Java server 启动后会先做端口就绪检查，再继续执行 Java/Go client 阶段。
- 默认 Java server 地址是 `127.0.0.1:20000`（可用环境变量覆盖）。

### 1.2 `start_integrate_test.sh` 与 `integrate_test.sh` 使用方式

前置条件：
- Docker / Docker Compose
- Go 工具链
- Maven（可选；如果要执行 Java 阶段则需要）

运行完整 CI 样例列表：
```bash
./start_integrate_test.sh
```

`start_integrate_test.sh` 会：
- 通过仓库根目录 `docker-compose.yml` 启动依赖
- 做依赖健康检查
- 逐个调用 `./integrate_test.sh <sample>`
- 在退出时（成功或失败）停止依赖

运行单个 sample：
```bash
./integrate_test.sh helloworld
./integrate_test.sh direct
```

常用环境变量：
- `GO_CLIENT_TIMEOUT_SECONDS`（默认：`90`）
- `JAVA_CLIENT_TIMEOUT_SECONDS`（默认：`180`）
- `JAVA_SERVER_READY_TIMEOUT_SECONDS`（默认：`60`）
- `JAVA_SERVER_HOST`（默认：`127.0.0.1`）
- `JAVA_SERVER_PORT`（默认：`20000`）

### 1.3 如何新增一个集成测试 sample

至少需要以下目录：
- `go-server/cmd/*.go`
- `go-client/cmd/*.go`

推荐结构：
- `go-server/conf/dubbogo.yml`（或 `dubbogo.yaml`）
- `go-client/conf/dubbogo.yml`（或 `dubbogo.yaml`）
- 可选 Java 互通：
  - `java-server/run.sh`
  - `java-client/run.sh`

Java 脚本要求：
- `java-server/run.sh` 需要保证 server 进程可持续运行（不要依赖后台 stdin）。
- `java-client/run.sh` 失败时必须返回非 0 退出码。

验证步骤：
1. 先单独跑 sample：
   ```bash
   ./integrate_test.sh <your-sample-path>
   ```
2. 将 sample 加入 `start_integrate_test.sh` 的 `array` 列表。
3. 再跑全量：
   ```bash
   ./start_integrate_test.sh
   ```
4. 确保失败路径可见（非 0 退出码、日志清晰）。

## 2. 手动运行（不使用脚本）

本节介绍如何不使用 `start_integrate_test.sh` 和 `integrate_test.sh`，手动运行单个 sample。

示例 sample：`helloworld`

### 2.1 启动依赖

```bash
cd <PATH OF dubbo-go-samples>
docker compose -f docker-compose.yml up -d
```

如果环境使用老版本 compose：
```bash
docker-compose -f docker-compose.yml up -d
```

### 2.2 运行 Go server

新开一个终端：
```bash
cd <PATH OF dubbo-go-samples>/helloworld
export DUBBO_GO_CONFIG_PATH=./go-server/conf/dubbogo.yml
go run ./go-server/cmd/*.go
```

### 2.3 运行 Go client

再开一个终端：
```bash
cd <PATH OF dubbo-go-samples>/helloworld
export DUBBO_GO_CONFIG_PATH=./go-client/conf/dubbogo.yml
go run ./go-client/cmd/*.go
```

### 2.4 可选：验证 Java 互通

1. 先停止 Go server（Go server 终端中 `Ctrl+C`）。
2. 启动 Java server：
   ```bash
   cd <PATH OF dubbo-go-samples>/helloworld/java-server
   bash ./run.sh
   ```
3. 运行 Java client：
   ```bash
   cd <PATH OF dubbo-go-samples>/helloworld/java-client
   bash ./run.sh
   ```
4. 再执行一次 Go client（2.3 的终端）验证 Go -> Java 调用。

### 2.5 清理

前台进程用 `Ctrl+C` 停止后，执行：
```bash
cd <PATH OF dubbo-go-samples>
docker compose -f docker-compose.yml down
```
