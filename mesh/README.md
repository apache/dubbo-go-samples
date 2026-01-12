# dubbo-go Proxy Mesh Using Istio

[English](README.md) | [中文](README_zh.md)

Follow the steps below to deploy the demo to your local cluster.

* [1 Overall Objectives](#target)
* [2 Basic Workflow](#basic)
* [3 Detailed Steps](#detail)
  + [3.1 Environment Requirements](#env)
  + [3.2 Prerequisites](#prepare)
  + [3.3 Building Images](#build)
  + [3.4 Deployment to Kubernetes](#deploy)
    - [3.4.1 Deploy Provider](#deploy_provider)
    - [3.4.2 Deploy Consumer](#deploy_consumer)
  + [3.5 Verify Provider and Consumer Communication](#check)
  + [3.6 Istio Traffic Governance](#traffic)

<h2 id="target">1 Overall Objectives</h2>

* Deploy dubbo-go applications to Kubernetes
* Enable automatic Envoy sidecar injection via Istio
* Implement traffic management based on Istio rules

<h2 id="basic">2 Basic Workflow & Architecture</h2>

This example demonstrates how to deploy dubbo-go applications within an Istio service mesh to achieve automatic Envoy proxying for Dubbo services.

Key steps to complete this example:

1. Create a dubbo-go application
2. Build container images
3. Deploy dubbo-go Provider and Consumer to Kubernetes and verify successful Envoy proxy injection
4. Verify that Envoy discovers service addresses, intercepts RPC traffic correctly, and implements load balancing
5. Apply traffic governance using Istio rules

<h2 id="detail">3 Detailed Steps</h2>

<h3 id="env">3.1 Environment Requirements</h3>

Ensure the following tools are installed locally to provide container runtime, Kubernetes cluster, and access tools:

* [Docker](https://www.docker.com/get-started/)
* [Minikube](https://minikube.sigs.k8s.io/docs/start/)
* [Kubectl](https://kubernetes.io/docs/tasks/tools/)
* [Istio](https://istio.io/latest/docs/setup/getting-started/)
* [Kubens(optional)](https://github.com/ahmetb/kubectx)

Start your local Kubernetes cluster:

```shell
minikube start
```

Verify the cluster is running properly and kubectl is bound to the default local cluster:

```shell
kubectl cluster-info
```

<h3 id="prepare">3.2 Prerequisites</h3>

Create a dedicated namespace `dubbo-demo` for the example project and enable automatic sidecar injection:

```shell
# Initialize namespace and enable sidecar injection
kubectl apply -f deploy/Namespace.yml

# Switch to the namespace
kubens dubbo-demo
```

<h3 id="build">3.3 Building Images</h3>

Build the Docker images using the provided script:

```shell
./build.sh
```

Or build manually:

```shell
# From the repository root
cd /path/to/dubbo-go-samples

# Build Provider image
docker build -f mesh/go-server/Dockerfile -t dubbo-go-mesh-provider:latest .

# Build Consumer image
docker build -f mesh/go-client/Dockerfile -t dubbo-go-mesh-consumer:latest .

# If using minikube, load images into minikube
minikube image load dubbo-go-mesh-provider:latest
minikube image load dubbo-go-mesh-consumer:latest
```

<h3 id="deploy">3.4 Deployment to Kubernetes</h3>

<h4 id="deploy_provider">3.4.1 Deploy Provider</h4>

```shell
# Deploy Service
kubectl apply -f deploy/provider/Service.yml

# Deploy Deployment
kubectl apply -f deploy/provider/Deployment.yml
```

The above commands create a Service named `server-demo`. Note that the service name matches the Dubbo application name defined in the code (via `dubbo.WithName("server-demo")`).

The Deployment creates a 2-replica pod instance. The Provider is now up and running.

You can check the startup logs with the following commands:

```shell
# Check pod list
kubectl get pods -l app=server-demo

# Check pod deployment logs
kubectl logs -l app=server-demo -c server
```

At this point, the pod should contain a dubbo provider container instance and an Envoy Sidecar container instance.

<h4 id="deploy_consumer">3.4.2 Deploy Consumer</h4>

```shell
# Deploy Service
kubectl apply -f deploy/consumer/Service.yml

# Deploy Deployment
kubectl apply -f deploy/consumer/Deployment.yml
```

Deploying the consumer follows the same process as the provider. Here we also keep the K8S Service name consistent with the Dubbo consumer application name (defined in the code via `dubbo.WithName("client-demo")`).

<h3 id="check">3.5 Verify Provider and Consumer Communication</h3>

After completing step 3.4, check the startup logs to verify that the consumer successfully consumes the provider service:

```shell
# Check pod list
kubectl get pods -l app=client-demo

# Check application logs
kubectl logs -l app=client-demo -c server

# Check istio-proxy logs
kubectl logs -l app=client-demo -c istio-proxy
```

For detailed log examples, refer to the [Java version demo](https://github.com/apache/dubbo-samples/tree/master/dubbo-samples-mesh-k8s)

<h3 id="traffic">3.6 Istio Traffic Governance</h3>

For traffic governance examples, refer to the [Java version demo](https://github.com/apache/dubbo-samples/tree/master/dubbo-samples-mesh-k8s)

#### View Dashboard

Check the Istio official documentation for [how to start the dashboard](https://istio.io/latest/docs/setup/getting-started/#dashboard).

### Server-Side Code Example

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

In a Mesh environment, the server does not require a registry center. It only needs to listen on a port (50052). Service discovery is handled by Kubernetes DNS and Istio.

### Client-Side Code Example

```go
// Use direct URL for mesh environment
cli, err := client.NewClient(
    client.WithClientURL("tri://server-demo.dubbo-demo.svc.cluster.local:50052"),
)

svc, err := greet.NewGreeter(cli)
reply, err := svc.SayHello(context.Background(), req)
```

In a Kubernetes mesh environment, the client connects directly using the service's DNS name in the format:

- `tri://service-name.namespace.svc.cluster.local:port`
- In this example: `tri://server-demo.dubbo-demo.svc.cluster.local:50052`
