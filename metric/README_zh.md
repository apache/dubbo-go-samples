# metric 示例

### 背景

[dubbo metric](https://github.com/alibaba/metrics) 是一套标准度量库,我们可以轻易列举出来在 RPC 领域里面我们所关心的各种指标，诸如每个服务的调用次数，响应时间；如果更加细致一点，还有各种响应时间的分布，平均响应时间.对原理感兴趣的可以去看下 [eBay 邓明：dubbo-go 中 metrics 的设计](https://dubbo.apache.org/zh/blog/2021/01/11/dubbo-go-中-metrics-的设计) 以下我将要使用[prometheus](https://prometheus.io/docs/introduction/overview) 来展示metric信息

### 效果图

![metric](../.images/metrics.png)

### 注意事项

* [prometheus](https://prometheus.io/docs/introduction/overview) 需要的是一个合法的名称而根据 [dubbo-go源码](https://github.com/apache/dubbo-go/blob/master/metrics/prometheus/reporter.go) 得知`namespace = config.GetApplicationConfig().Name`也就是说你服务的名称字母、数组、下划线才可以

* 需要修改[prometheus.yml](./go-server/docker/config/prometheus.yml)ip为本地物理机器的ip

> prometheus.yml配置如下
```yaml
# my global config
global:
  scrape_interval: 120s
  evaluation_interval: 120s
  external_labels:
    monitor: 'metric-dubbo-go-server'
scrape_configs:
  - job_name: 'prometheus'
    scrape_interval: 120s
    static_configs:
      - targets: [ 'localhost:9090' ]

  - job_name: 'metric-dubbo-server'
    scheme: http
    scrape_interval: 10s
    static_configs:
      # 这里需要写本机机器的ip
      - targets: [ '本地ip:8080' ]
```

请参阅根目录中的 [HOWTO.md](../HOWTO_zh.md) 来运行本例。
