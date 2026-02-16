
# Dubbo-Go Kubernetes Probe Sample

English | [中文](README_CN.md)

This example demonstrates how to integrate **Kubernetes liveness, readiness, and startup probes** with a Dubbo-Go Triple service.

It showcases how to expose dedicated probe endpoints and how Kubernetes reacts during application startup and warm-up.


## Project Layout

```
metrics/probe/
├── go-server/
│   ├── cmd/main.go        # Application entrypoint
│   ├── build.sh           # Docker build script
│   └── Dockerfile         # Container image definition
└── deploy/
    └── server-deployment.yml  # Kubernetes deployment
```


## Default Configuration

| Item                | Value                         |
| ------------------- | ----------------------------- |
| Triple service port | `20000`                       |
| Probe HTTP port     | `22222`                       |
| Probe endpoints     | `/live`, `/ready`, `/startup` |
| Client target       | `tri://127.0.0.1:20000`       |

### Probe Semantics

* `GET /live`
  Indicates whether the process is alive.

* `GET /ready`
  Indicates whether the service is ready to receive traffic.

* `GET /startup`
  Indicates whether the application startup phase has completed.


## Run Locally

### 1️Start the server

```bash
go run ./metrics/probe/go-server/cmd/main.go
```

### 2️Monitor probe endpoints

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

### Expected Behavior

| Phase                | /live | /ready | /startup |
| -------------------- | ----- | ------ | -------- |
| Process started      | 200   | 503    | 503      |
| During warm-up       | 200   | 503    | 503      |
| After warm-up (~15s) | 200   | 200    | 200      |


## Kubernetes Probe Configuration

Example probe configuration:

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


## Run on Kubernetes

### Build the image

From the repository root:

```bash
./metrics/probe/go-server/build.sh
```


### Load image into local cluster (if needed)

For example, with Minikube:

```bash
minikube image load dubbo-go-probe-server:latest
```


### Deploy to Kubernetes

```bash
kubectl apply -f metrics/probe/deploy/server-deployment.yml
kubectl rollout status deploy/dubbo-go-probe-server
kubectl get pod -l app=dubbo-go-probe-server
```


### Inspect probe status

```bash
kubectl describe pod -l app=dubbo-go-probe-server
```

### Expected Kubernetes Behavior

* Immediately after deployment:

  * `Ready` = `False`
  * `ContainersReady` = `False`

* After ~15 seconds (warm-up completed):

  * `Ready` = `True`
  * `ContainersReady` = `True`

This demonstrates proper separation of:

* Process liveness
* Service readiness
* Startup lifecycle completion