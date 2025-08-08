# Dubbo-Go Metrics Monitoring Example

English | [中文](README_CN.md)

This example demonstrates how to use the Push and Pull modes of **Prometheus Pushgateway** to monitor a Dubbo-Go application and visualize the data with Grafana.

-----

## Core Architecture

The monitoring data flow is as follows:

**Push Mode: Application (go-client / go-server) -> Prometheus Pushgateway -> Prometheus -> Grafana**

**Pull Mode: Application (go-client / go-server) -> Prometheus -> Grafana**

## Included Components

| Component | Port | Description |
| :--- | :--- | :--- |
| **Grafana** | `3000` | A dashboard for visualizing metrics. |
| **Prometheus** | `9090` | Responsible for storing and querying metric data. It pulls data from the Pushgateway. |
| **Pushgateway** | `9091` | Used to receive metrics pushed from the Dubbo-Go application. |
| **go-server** | N/A | Dubbo-Go service provider (Provider) example. |
| **go-client** | N/A | Dubbo-Go service consumer (Consumer) example that continuously calls the server. |

## 🚀 Quick Start

Please follow the steps below to run this example.

### Prerequisites:

- Please configure the network addresses in `prometheus_pull.yml`, `prometheus_push.yml`, `go-client/cmd/main.go`, and `go-server/cmd/main.go` according to your actual network environment.
- If you want to try the Push mode, change line 38 in `docker-compose.yml` from `- ./prometheus_pull.yml:/etc/prometheus/prometheus.yml` to `- ./prometheus_push.yml:/etc/prometheus/prometheus.yml`, and then restart the services.

### Step 1: Start the Monitoring Stack

First, start the Grafana, Prometheus, and Pushgateway services. We use `docker-compose` to do this with a single command.

```bash
# Enter the metrics directory
cd metrics
# Start all monitoring services in the background
docker-compose up -d
```

You can now access the web UI for each service at the following addresses:

- **Grafana**: `http://localhost:3000`
- **Prometheus**: `http://localhost:9090`
- **Pushgateway**: `http://localhost:9091`

### Step 2: Start the Dubbo-Go Server

In the metrics directory, open a new terminal window and run the server program.

```bash
go run ./go-server/cmd/main.go
```

You will see logs indicating that the server has started successfully and registered its services.

### Step 3: Start the Dubbo-Go Client

In the metrics directory, open another new terminal window and run the client program. The client will continuously call the server's methods, with random failures to generate monitoring metrics.

```bash
go run ./go-client/cmd/main.go
```

The client will start printing call results while pushing monitoring metrics to the Pushgateway. You can see the pushed metrics on the Pushgateway UI (`http://localhost:9091/metrics`).

### Step 4: Configure Grafana and Import the Dashboard

Now that all services are running, let's configure Grafana to display the data.

#### 4.1. Add Prometheus Data Source

1.  Open the Grafana website: [`http://localhost:3000`](https://www.google.com/search?q=http://localhost:3000) (default username/password: `admin`/`admin`).
2.  In the left-side menu, navigate to **Home -> Connections -> Data sources**.
3.  Click the **[Add new data source]** button.
4.  Select **Prometheus**.
5.  In the **Prometheus server URL** field, enter `http://host.docker.internal:9090`.
    > **Note**: `host.docker.internal` is a special DNS name that allows Docker containers (like Grafana) to access the host machine's network. You can configure it according to your actual situation.
6.  Click the **[Save & test]** button at the bottom. You should see a "Data source is working" success message.

#### 4.2. Import the Dubbo Monitoring Dashboard

1.  In the left-side menu, navigate to **Home -> Dashboards**.
2.  Click **[New]** -> **[Import]** in the top right corner.
3.  Copy the contents of `grafana.json` into the **Import via panel json** text box, or click the **Upload JSON file** button to upload the `grafana.json` file.
4.  On the next page, make sure to select the Prometheus data source we just created for the dashboard.
5.  Click the **[Import]** button.

### Step 5: View the Monitoring Dashboard

After a successful import, you will see a complete Dubbo observability dashboard\! The data in the panels (like QPS, success rate, latency, etc.) will update dynamically as the client continues to make calls.

Enjoy\!

## Troubleshooting

- **Grafana dashboard shows "No Data"**

   - Verify that the Prometheus data source URL (`http://host.docker.internal:9090`) is correct and that the connection test was successful.
   - Go to the Prometheus UI (`http://localhost:9090`), and check the `Status -> Targets` page to ensure the `pushgateway` job has a status of **UP**.
   - In the Prometheus query bar, enter `dubbo_consumer_requests_succeed_total` to confirm that data can be queried.

- **Cannot connect to `host.docker.internal`**

   - `host.docker.internal` is a built-in feature of Docker. If this address is not accessible, replace the IP address in `metrics/prometheus.yml` and the Grafana data source address with your actual IP address.

-----

## Deploying to Kubernetes

#### kube-prometheus

To install Prometheus in Kubernetes (k8s), please refer to the [kube-prometheus](https://github.com/prometheus-operator/kube-prometheus) project.

Set the service type in `prometheus-service.yaml` to `NodePort`.

1.  Add the `dubboPodMoitor.yaml` file to the `manifests` directory of `kube-prometheus` with the following content:

    ```yaml
    apiVersion: monitoring.coreos.com/v1
    kind: PodMonitor
    metadata:
      name: podmonitor
      labels:
        app: podmonitor
      namespace: monitoring
    spec:
      namespaceSelector:
        matchNames:
          - dubbo-system
      selector:
        matchLabels:
          app-type: dubbo
      podMetricsEndpoints:
        - port: metrics # Reference the port name 'metrics' of the dubbo-app
          path: /prometheus
    ---
    # Role-Based Access Control (RBAC)
    apiVersion: rbac.authorization.k8s.io/v1
    kind: Role
    metadata:
      namespace: dubbo-system
      name: pod-reader
    rules:
      - apiGroups: [ "" ]
        resources: [ "pods" ]
        verbs: [ "get", "list", "watch" ]

    ---
    # Role-Based Access Control (RBAC)
    apiVersion: rbac.authorization.k8s.io/v1
    kind: RoleBinding
    metadata:
      name: pod-reader-binding
      namespace: dubbo-system
    roleRef:
      apiGroup: rbac.authorization.k8s.io
      kind: Role
      name: pod-reader
    subjects:
      - kind: ServiceAccount
        name: prometheus-k8s
        namespace: monitoring
    ```

2.  Execute `kubectl apply -f Deployment.yaml`

3.  Open the Prometheus web interface, for example, `http://localhost:9090/targets`