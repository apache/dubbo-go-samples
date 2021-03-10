# Sentinel Filter 示例

### 背景

在 dubbo-go 中，[sentinel](https://github.com/alibaba/sentinel) 是作为 filter 的形式集成的。这个示例展示了如何在 dubbo-go 中配置并使用 sentinel。

### 示例

##### 1. 代码

通过代码初始化 sentinel 并配置相关的限流规则，如下所示：

```go
var err error
conf := sentinelConf.NewDefaultConfig()
err = api.InitWithConfig(conf)

_, err = flow.LoadRules([]*flow.Rule{
    {
        // protocol:consumer:interfaceName:group:version:method
        Resource:               "dubbo:consumer:org.apache.dubbo.UserProvider:::GetUser()",
        TokenCalculateStrategy: flow.Direct,
        ControlBehavior:        flow.Reject,
        Threshold:              1,
        StatIntervalInMs:       1000,
    },
})
```

##### 2. 配置

Sentinel filter 分为 sentinel-consumer 和 sentinel-provider 两种，分别用于 consumer 端和 provider 端，如下所示：

```yaml
# consumer filter config
filter: "sentinel-consumer"

# provider filter config
filter: "sentinel-provider"
```

##### 3. 运行

请参阅根目录中的 [HOWTO.md](../../HOWTO_zh.md) 来运行本例。

```bash
[2021-03-10/10:55:46 main.main: client.go: 80] error: SentinelBlockError: FlowControl, message: flow reject check blocked
```
