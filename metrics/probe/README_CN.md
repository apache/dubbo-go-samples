
# Dubbo-Go Kubernetes 探针示例

[English](README.md) | 中文


本示例演示如何在 Dubbo-Go Triple 服务中集成 **Kubernetes 的 liveness、readiness 和 startup 探针**。

通过暴露独立的探针端点，可以清晰地观察应用启动、预热以及就绪阶段在 Kubernetes 中的行为。


## 项目结构

```
metrics/probe/
├── go-server/
│   ├── cmd/main.go              # 程序入口
│   ├── build.sh                 # Docker 构建脚本
│   └── Dockerfile               # 镜像定义文件
└── deploy/
    └── server-deployment.yml    # Kubernetes 部署文件
```


## 默认配置

| 项目          | 默认值                           |
| ----------- | ----------------------------- |
| Triple 服务端口 | `20000`                       |
| 探针 HTTP 端口  | `22222`                       |
| 探针路径        | `/live`, `/ready`, `/startup` |
| 客户端目标地址     | `tri://127.0.0.1:20000`       |


## 探针语义说明

* `GET /live`
  表示进程是否存活（进程级健康检查）。

* `GET /ready`
  表示服务是否已准备好接收流量。

* `GET /startup`
  表示应用是否完成启动阶段。


## 本地运行

### 启动服务

```bash
go run ./metrics/probe/go-server/cmd/main.go
```


### 观察探针状态

```bash
watch -n 1 '
for p in live ready startup; do
  url="http://127.0.0.1:22222/$p"

  body=$(curl -sS --max-time 2 "$url" 2>&1)
  code=$(curl -s -o /dev/null --max-time 2 -w "%{http_code}" "$url" 2>/dev/null)

  printf "%-8s [%s] %s\n" "$p" "$code" "$body"
done
'
```


### 预期行为

| 阶段           | /live | /ready | /startup |
| ------------ | ----- | ------ | -------- |
| 进程刚启动        | 200   | 503    | 503      |
| 预热阶段         | 200   | 503    | 503      |
| 预热完成（约 15 秒） | 200   | 200    | 200      |

说明：

* `/live` 只要进程未崩溃就返回 200。
* `/ready` 与 `/startup` 在应用完成预热前返回 503。
* 预热完成后，三个端点均返回 200。


## Kubernetes 探针配置示例

```yaml
livenessProbe:
  httpGet:
    path: /live
    port: 22222

readinessProbe:
  httpGet:
    path: /ready
    port: 22222

startupProbe:
  httpGet:
    path: /startup
    port: 22222
```


## 在 Kubernetes 中运行

### 构建镜像

在仓库根目录执行：

```bash
./metrics/probe/go-server/build.sh
```


### 将镜像加载到本地集群

例如使用 Minikube：

```bash
minikube image load dubbo-go-probe-server:latest
```


### 部署到 Kubernetes

```bash
kubectl apply -f metrics/probe/deploy/server-deployment.yml
kubectl rollout status deploy/dubbo-go-probe-server
kubectl get pod -l app=dubbo-go-probe-server
```


### 查看探针状态

```bash
kubectl describe pod -l app=dubbo-go-probe-server
```


### 预期 Kubernetes 行为

* 刚部署时：

  * `Ready` = `False`
  * `ContainersReady` = `False`

* 约 15 秒后（预热完成）：

  * `Ready` = `True`
  * `ContainersReady` = `True`

这体现了：

* 进程存活检查（Liveness）
* 服务可用性检查（Readiness）
* 启动阶段控制（Startup）

三者的职责分离。
