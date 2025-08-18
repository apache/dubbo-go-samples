# Dubbo-Go Pushgateway Cleaner Documentation

## Overview

Dubbo-Go Pushgateway Cleaner is a tool designed to solve the "zombie metrics" issue in Prometheus Pushgateway. It provides an operation-side cleanup solution.  
(This solution is based on the `job_pushed_at_seconds` metric in `Prometheus Pushgateway`.)

---

### Features

-   **TTL Cleanup**: Delete metrics older than a specified time.
-   **Job Filtering**: Filter by job name prefix.
-   **Dual Running Modes**: Supports scheduled tasks and daemon mode.

### Configuration Parameters

| Parameter       | Short | Default  | Description                      |
|-----------------|-------|----------|----------------------------------|
| `-pushgw`       | -     | Required | Pushgateway address              |
| `-ttl`          | -     | `3600`   | Metric retention time (seconds)  |
| `-job-prefix`   | -     | `*`      | Job name prefix filter           |
| `-daemon`       | -     | `false`  | Enable daemon mode               |
| `-interval`     | -     | `60`     | Cleanup interval (seconds)       |

### Running Modes

#### Scheduled Task Mode (Recommended)

```bash
./pgw-cleaner \
  -pushgw http://pushgateway:9091 \
  -ttl 3600 \
  -job-prefix dubbo_
```
-   Performs a single cleanup and exits
-   Suitable for running as a CronJob

#### Daemon Mode

```bash
./pgw-cleaner \
-pushgw http://pushgateway:9091 \
-ttl 3600 \
-job-prefix dubbo_ \
-daemon \
-interval 300
```

-   Runs continuously with periodic cleanup
-   Exposes metrics at: `http://<ip>:9105/metrics`

* * *

## 3. Containerized Deployment

### Build Image

```
docker build . -t pgw-cleaner:latest
```

### Docker Run

#### Single Execution Mode

```bash
docker run --rm pgw-cleaner:latest \
-pushgw http://127.0.0.1:9091 \
-job-prefix dubbo_ \
-ttl 3600
```

#### Daemon Mode

```bash
docker run -d -p 9105:9105 pgw-cleaner:latest \
-pushgw http://127.0.0.1:9091 \
-job-prefix dubbo_ \
-ttl 3600 \
-daemon \
-interval 300
```

* * *

## 4. Kubernetes Deployment

### Prerequisites

-   A running Kubernetes cluster (minikube / k3s / cloud-based).

-   `kubectl` configured to point to the cluster.

-   Pushgateway deployed inside the cluster and accessible via DNS, e.g.:

    ```
    http://pushgateway.monitoring.svc.cluster.local:9091
    ```

-   Cleaner image built and pushed to an image registry, e.g.:

    ```
    pgw-cleaner:latest
    ```

-   Ensure the `monitoring` namespace exists. If not, create it:

```bash
kubectl create namespace monitoring
```

### CronJob Mode (Recommended)

```bash
kubectl apply -f ./deploy/cronjob.yaml
```

**Advantages**:

-   Low resource usage
-   Executes on schedule
-   Automatic retries on failure

### Deployment Mode

```bash
kubectl apply -f ./deploy/deployment.yaml
```

**Advantages**:

-   Continuous metric monitoring
-   Exposes its own monitoring metrics
-   Suitable for frequent cleanup scenarios

* * *

## 5. Best Practices

### TTL Recommendation

| Scenario   | Recommended TTL | Notes                                  |
| ---------- | --------------- | -------------------------------------- |
| Short jobs | 5-10 minutes    | Short-lived tasks, quick cleanup       |
| Daemons    | 1-2 hours       | Allows for brief service interruptions |
| Production | 30 minutes      | Balance between freshness and cost     |

### Job Prefix Convention

**Recommended format**: `<env>-<service>-`  
Examples:

-   `prod-order-service-`
-   `dev-user-service-`

**Benefits**:

1.  Prevents accidental deletion of other metrics
1.  Supports environment/service-based cleanup
1.  Improves maintainability

### Monitoring Metrics

Metrics exposed in daemon mode:

-   `pgw_cleaner_checked_total`: Total number of jobs checked
-   `pgw_cleaner_deleted_total`: Total number of jobs deleted
-   `pgw_cleaner_errors_total`: Total cleanup errors

Grafana example query:

```
sum(rate(pgw_cleaner_deleted_total[5m])) by (job)
```

* * *

## 6. Troubleshooting

### 1. Cleanup Ineffective

**Issue**: Metrics not cleaned, expired data still present in Pushgateway  
**Steps**:

1.  Confirm `job_pushed_at_seconds` exists:

```bash
curl http://pushgateway:9091/metrics | grep job_pushed_at_seconds
```

2.  Check TTL validity:

    -   Current timestamp: `date +%s`
    -   Metric timestamp: obtained from Pushgateway
    -   Ensure `(current_time - metric_time) > TTL`

2.  Verify job prefix match:

```bash
# List all jobs
curl http://pushgateway:9091/api/v1/metrics
```

### 2. Missing Timestamp Metric

**Issue**: `job_pushed_at_seconds` not reported  
**Solution**:

1.  **Dubbo-Go Client Config**:

    ``` go
    pushedAt := prometheus.NewGauge(prometheus.GaugeOpts{
        Name: "job_pushed_at_seconds",
        Help: "Last push timestamp",
    })
    prometheus.MustRegister(pushedAt)

    // Update before each push
    func pushMetrics() {
        pushedAt.SetToCurrentTime()
        // ... push logic ...
    }
    ```

1.  **Scheduled Verification**:

```bash
./pgw-cleaner -verify -pushgw http://pushgateway:9091
```

### 3. Authentication Failure

**Issue**: `401 Unauthorized`  
**Solution**:

```bash
# Add authentication in the cleanup command
./pgw-cleaner \
-pushgw http://username:password@pushgateway:9091 \
-ttl 1h
```

### 4. Performance Issues

**Issue**: Cleanup takes too long  
**Optimization**:

1.  Increase concurrency:

```bash
./pgw-cleaner -workers=5
```

2.  Split cleanup:

```
# Split by environment
- -job-prefix=prod-east-
- -job-prefix=prod-west-
```

3.  Adjust time window:

```bash
# Cleanup within a 1-hour window
./pgw-cleaner -time-window=1h
```