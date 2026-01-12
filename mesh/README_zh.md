# dubbo-go Proxy Mesh Using Istio

[English](README.md) | [中文](README_zh.md)

可以按照下文步骤，将 Demo 部署到本地集群。

* [1 总体目标](#target)
* [2 基本流程](#basic)
* [3 详细步骤](#detail)
  + [3.1 环境要求](#env)
  + [3.2 前置条件](#prepare)
  + [3.3 构建镜像](#build)
  + [3.4 部署到 Kubernetes](#deploy)
    - [3.4.1 部署 Provider](#deploy_provider)
    - [3.4.2 部署 Consumer](#deploy_consumer)
  + [3.5 检查 Provider 和 Consumer 正常通信](#check)
  + [3.6 Istio 流量治理](#traffic)

<h2 id="target">1 总体目标</h2>

* 部署 dubbo-go 应用到 Kubernetes
* Istio 自动注入 Envoy 并实现流量拦截
* 基于 Istio 规则进行流量治理

<h2 id="basic">2 基本流程与工作原理</h2>

这个示例演示了如何将 dubbo-go 开发的应用部署在 Istio 体系下，以实现 Envoy 对 Dubbo 服务的自动代理。

完成示例将需要的步骤如下：

1. 创建一个 dubbo-go 应用
2. 构建容器镜像
3. 分别部署 dubbo-go Provider 与 dubbo-go Consumer 到 Kubernetes 并验证 Envoy 代理注入成功
4. 验证 Envoy 发现服务地址、正常拦截 RPC 流量并实现负载均衡
5. 基于 Istio 规则进行流量治理

<h2 id="detail">3 详细步骤</h2>

<h3 id="env">3.1 环境要求</h3>

请确保本地安装如下环境，以提供容器运行时、Kubernetes集群及访问工具

* [Docker](https://www.docker.com/get-started/)
* [Minikube](https://minikube.sigs.k8s.io/docs/start/)
* [Kubectl](https://kubernetes.io/docs/tasks/tools/)
* [Istio](https://istio.io/latest/docs/setup/getting-started/)
* [Kubens(optional)](https://github.com/ahmetb/kubectx)

通过以下命令启动本地 Kubernetes 集群

```shell
minikube start
```

通过 kubectl 检查集群正常运行，且 kubectl 绑定到默认本地集群

```shell
kubectl cluster-info
```

<h3 id="prepare">3.2 前置条件</h3>

通过以下命令为示例项目创建独立的 Namespace `dubbo-demo`，同时开启 sidecar 自动注入。

```shell
# 初始化命名空间并开启sidecar自动注入
kubectl apply -f deploy/Namespace.yml

# 切换命名空间
kubens dubbo-demo
```

<h3 id="build">3.3 构建镜像</h3>

使用提供的构建脚本构建 Docker 镜像：

```shell
./build.sh
```

或者手动构建：

```shell
# 从仓库根目录执行
cd /path/to/dubbo-go-samples

# 构建服务端镜像
docker build -f mesh/go-server/Dockerfile -t dubbo-go-mesh-provider:latest .

# 构建客户端镜像
docker build -f mesh/go-client/Dockerfile -t dubbo-go-mesh-consumer:latest .

# 如果使用 minikube，加载镜像到 minikube
minikube image load dubbo-go-mesh-provider:latest
minikube image load dubbo-go-mesh-consumer:latest
```

<h3 id="deploy">3.4 部署到 Kubernetes</h3>

<h4 id="deploy_provider">3.4.1 部署 Provider</h4>

```shell
# 部署 Service
kubectl apply -f deploy/provider/Service.yml

# 部署 Deployment
kubectl apply -f deploy/provider/Deployment.yml
```

以上命令创建了一个名为 `server-demo` 的 Service，注意这里的 service name 与项目中的 dubbo 应用名是一样的（在代码中通过 `dubbo.WithName("server-demo")` 定义）。

接着 Deployment 部署了一个 2 副本的 pod 实例，至此 Provider 启动完成。

可以通过如下命令检查启动日志。

```shell
# 查看 pod 列表
kubectl get pods -l app=server-demo

# 查看 pod 部署日志
kubectl logs -l app=server-demo -c server
```

这时 pod 中应该有一个 dubbo provider 容器实例，同时还有一个 Envoy Sidecar 容器实例。

<h4 id="deploy_consumer">3.4.2 部署 Consumer</h4>

```shell
# 部署 Service
kubectl apply -f deploy/consumer/Service.yml

# 部署 Deployment
kubectl apply -f deploy/consumer/Deployment.yml
```

部署 consumer 与 provider 是一样的，这里也保持了 K8S Service 与 Dubbo consumer application name（在代码中通过 `dubbo.WithName("client-demo")` 定义）一致。

<h3 id="check">3.5 检查 Provider 和 Consumer 正常通信</h3>

继执行 3.4 步骤后，检查启动日志，查看 consumer 完成对 provider 服务的消费。

```shell
# 查看 pod 列表
kubectl get pods -l app=client-demo

# 查看 pod 部署日志
kubectl logs -l app=client-demo -c server

# 查看 pod istio-proxy 日志
kubectl logs -l app=client-demo -c istio-proxy
```

具体日志情况可参考 [Java 版本对应 demo](https://github.com/apache/dubbo-samples/tree/master/dubbo-samples-mesh-k8s)

<h3 id="traffic">3.6 Istio 流量治理</h3>

参考 [Java 版本对应 demo](https://github.com/apache/dubbo-samples/tree/master/dubbo-samples-mesh-k8s)

#### 查看 dashboard

Istio 官网查看 [如何启动 dashboard](https://istio.io/latest/docs/setup/getting-started/#dashboard)。

### 服务端代码示例

```go
// In mesh environment, server only needs to listen on port
// Service discovery is handled by Kubernetes DNS and Istio
srv, err := server.NewServer(
    server.WithServerProtocol(
        protocol.WithTriple(),
        protocol.WithPort(50052),
    ),
)

greet.RegisterGreeterHandler(srv, &GreeterProvider{})
srv.Serve()
```

在 Mesh 环境中，服务端无需注册中心，只需监听端口（50052）。服务发现由 Kubernetes DNS 和 Istio 负责。

### 客户端代码示例

```go
// Use direct URL for mesh environment
cli, err := client.NewClient(
    client.WithClientURL("tri://server-demo.dubbo-demo.svc.cluster.local:50052"),
)

svc, err := greet.NewGreeter(cli)
reply, err := svc.SayHello(context.Background(), req)
```

在 Kubernetes mesh 环境中，客户端直接使用服务的 DNS 名称连接，格式为：

- `tri://服务名.命名空间.svc.cluster.local:端口`
- 本示例中为：`tri://server-demo.dubbo-demo.svc.cluster.local:50052`