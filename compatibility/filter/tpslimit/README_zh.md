# TPSLimit Filter 示例

### 背景

Dubbo-go 内置了限流 filter "tpslimit"。可以通过在服务端的配置来激活，另外，用户还可以自定义限流策略和拒绝访问后的处理逻辑。

### 示例

##### 1. 代码

A) 自定义限流策略：

通过实现 filter.TpsLimitStrategy 来自定义限流策略。在本例中，采取的策略是随机限流。例子链接为：https://github.com/apache/dubbo-go-samples/tree/master/filter/tpslimit/go-server/pkg/limit_strategy.go

B) 自定义拒绝访问处理：

通过实现 filter.RejectedExecutionHandler。在本例中，当限流条件满足的情况下，拒绝访问的自定义处理逻辑将会返回 "The request is rejected and doesn't have any default value. " 的错误给客户端。例子链接为：https://github.com/apache/dubbo-go-samples/tree/master/filter/tpslimit/go-server/pkg/reject_handler.go

##### 2. 配置

在服务端的配置文件中，按如下所示配置该 filter：

```yaml
# service config
services:
  UserProvider:
    registry: demoZk
    protocol: dubbo
    interface: org.apache.dubbo.UserProvider
    tps.limiter: method-service
    tps.limit.strategy: RandomLimitStrategy
    tps.limit.rejected.handler: DefaultValueHandler
    tps.limit.interval: 5000
    tps.limit.rate: 300
```

##### 3. 运行

请参阅根目录中的 [HOWTO.md](../../../HOWTO_zh.md) 来运行本例。

观察服务端的输出：

```bash
[2021-03-10/17:11:10 github.com/apache/dubbo-go-samples/filter/tpslimit/go-server/pkg.RandomTpsLimitStrategy.IsAllowable: limit_strategy.go: 56] %s
Random IsAllowable!
2021-03-10T17:11:10.748+0800 ERROR   filter_impl/tps_limit_filter.go:69      The invocation was rejected due to over the tps limitation, ...
```

观察客户端的输出：

```bash
error: The request is rejected and doesn't have any default value. 
```