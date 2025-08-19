# Dubbo-Go Pushgateway Cleaner 文档

[English](README.md) | 中文

## 概述

Dubbo-Go Pushgateway Cleaner 是专门为解决 Prometheus Pushgateway 中"僵尸指标"问题设计的工具。现提供运维端清理方案。
(本方案基于`Prometheus Pushgateway`中的 `job_pushed_at_seconds` 指标进行实现)

* * *

### 功能特性

-   **TTL 清理**：删除超过指定时间的指标
-   **任务过滤**：按任务名前缀过滤
-   **双模式运行**：支持定时任务和常驻服务

### 配置参数

| 参数            | 缩写 | 默认值     | 说明             |
| ------------- | -- |---------|----------------|
| `-pushgw`     | -  | 必填      | Pushgateway 地址 |
| `-ttl`        | -  | `3600`  | 指标保留时间 (秒)     |
| `-job-prefix` | -  | `*`     | 任务名前缀过滤        |
| `-daemon`     | -  | `false` | 启用守护进程模式       |
| `-interval`   | -  | `60`    | 清理间隔(秒)        |

### 运行模式

#### 定时任务模式 (推荐)

```bash
./pgw-cleaner \
  -pushgw http://pushgateway:9091 \
  -ttl 3600 \
  -job-prefix dubbo_
```

-   执行单次清理后退出
-   适合通过 CronJob 定期执行

#### 守护进程模式

```bash
./pgw-cleaner \
  -pushgw http://pushgateway:9091 \
  -ttl 3600 \
  -job-prefix dubbo_ \
  -daemon \
  -interval 300
```
-   持续运行并定期清理
-   暴露监控指标：`http://<ip>:9105/metrics`

* * *

## 3. 容器化部署



### 构建镜像

```bash
docker build . -t pgw-cleaner:latest
```

### Docker 运行

####  单次执行模式
```bash
docker run --rm pgw-cleaner:latest \
  -pushgw http://127.0.0.1:9091 \
  -job-prefix dubbo_ \
  -ttl 3600
```
####  守护进程模式
```bash
docker run -d -p 9105:9105 pgw-cleaner:latest \
  -pushgw http://127.0.0.1:9091 \
  -job-prefix dubbo_ \
  -ttl 3600 \
  -daemon \
  -interval 300
```


* * *

## 4. Kubernetes 部署

### 前置条件
-   已有 Kubernetes 集群（minikube / k3s / 云上集群均可）。

-   `kubectl` 已配置指向集群。

-   Pushgateway 已部署在集群内，并可通过 DNS 访问，例如：

    ```
    http://pushgateway.monitoring.svc.cluster.local:9091
    ```

-   Cleaner 镜像已经构建并推送到镜像仓库，例如：

    ```
    pgw-cleaner:latest
    ```
- 确保 monitoring 命名空间存在，如果不存在可以创建：

```bash
kubectl create namespace monitoring
   ```

### CronJob 模式 (推荐)
```bash
kubectl apply -f ./deploy/cronjob.yaml
  ```

**特点**：

-   资源消耗低
-   按计划定期执行
-   失败自动重试

### Deployment 模式
```bash
kubectl apply -f ./deploy/deployment.yaml
```

**特点**：

-   持续监控指标
-   暴露自身监控数据
-   适合高频清理场景

* * *

## 5. 最佳实践

### TTL 设置建议

| 场景   | 推荐 TTL  | 说明           |
| ---- | ------- | ------------ |
| 短时任务 | 5-10 分钟 | 任务执行时间短，快速清理 |
| 常驻服务 | 1-2 小时  | 允许服务短暂中断     |
| 生产环境 | 30 分钟   | 平衡实时性与资源消耗   |

### 任务名前缀规范

**推荐格式**：`<环境>-<服务名>-`  
示例：
-   `prod-order-service-`
-   `dev-user-service-`

**优势**：

1.  避免误删其他服务指标
1.  支持按环境/服务清理
1.  提高可维护性

### 监控指标
守护进程模式暴露的指标：

-   `pgw_cleaner_checked_total`：已检查任务数
-   `pgw_cleaner_deleted_total`：已删除任务数
-   `pgw_cleaner_errors_total`：清理错误数

Grafana 监控示例：
```
sum(rate(pgw_cleaner_deleted_total[5m])) by (job)
```

* * *

## 6. 常见问题解决方案

### 1. 清理未生效

**现象**：指标未被清理，Pushgateway 中仍有过期数据  
**排查步骤**：

1.  确认`job_pushed_at_seconds`指标存在：
```bash
    curl http://pushgateway:9091/metrics | grep job_pushed_at_seconds
  ```

2. 检查 TTL 设置是否合理：

    -   当前时间戳：`date +%s`
    -   指标时间戳：从 Pushgateway 获取
    -   确保：`(当前时间 - 指标时间) > TTL`

3. 验证任务名前缀匹配：

 ```bash
    # 查看所有任务
    curl http://pushgateway:9091/api/v1/metrics
  ```

### 2. 时间戳指标缺失

**现象**：`job_pushed_at_seconds`指标未上报  
**解决方案**：

1.  **Dubbo-Go 客户端配置**：

    ``` go
    pushedAt := prometheus.NewGauge(prometheus.GaugeOpts{
        Name: "job_pushed_at_seconds",
        Help: "Last push timestamp",
    })
    prometheus.MustRegister(pushedAt)

    // 每次推送时更新
    func pushMetrics() {
        pushedAt.SetToCurrentTime()
        // ... 推送逻辑 ...
    }
    ```

2. **定时任务中检查**：
```bash
    ./pgw-cleaner -verify -pushgw http://pushgateway:9091
  ```

### 3. 认证失败

**现象**：`401 Unauthorized`错误  
**解决方案**：

```bash
# 在清理命令中添加认证
./pgw-cleaner \
  -pushgw http://username:password@pushgateway:9091 \
  -ttl 1h
```

### 4. 性能问题

**现象**：清理耗时过长  
**优化方案**：

1.  增加并行度：
```bash
    ./pgw-cleaner -workers=5
  ```
2. 分片清理：
```yaml
    # 按环境分片
    - -job-prefix=prod-east-
    - -job-prefix=prod-west-
   ```

3. 调整时间窗口：

 ```bash
    # 每次只清理1小时时间窗口
    ./pgw-cleaner -time-window=1h
   ```