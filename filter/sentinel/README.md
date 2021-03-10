# Sentinel Filter Sample

### Background

In dubbo-go, [sentinel](https://github.com/alibaba/sentinel) is supported as a filter. This example demonstrates how sentinel is configured and used in dubbo-go. 

### Example

##### 1. Code

Use the following code to initialize sentinel and the corresponding flow rule:

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

##### 2. Configuration

There are two builtin sentinel filters: sentinel-consumer and sentinel-provider. The former filter is used on consumer side, and the latter is for provider side. Check the following snippet for how to config it:

```yaml
# consumer filter config
filter: "sentinel-consumer"

# provider filter config
filter: "sentinel-provider"
```

##### 3. Run

Pls. refer to [HOWTO.md](../../HOWTO.md) under the root directory to run this sample.

```bash
[2021-03-10/10:55:46 main.main: client.go: 80] error: SentinelBlockError: FlowControl, message: flow reject check blocked
```
