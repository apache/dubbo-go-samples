## 日志使用

 该 samples 演示了如何使用 lumberjack 配置 dubbo-go logger

### 目录

* default:  默认打印到控制台
* level:    设置日志的隔离级别
* rolling:  输出到文件
* custom: 自定义 logger
* trace-integration: 集成 OpenTelemetry trace 信息

#### 默认配置

在配置文件中不添加 logger 配置日志将会打印到控制台, 也可在配置文件中配置日志, 可参照如下方式: 

zap 日志格式和级别设置
```yaml
    logger:
      driver: zap
      level: info
      format: text
      appender: console
      file:
        name: logs.log
        max-size: 1
        max-age: 3
        max-backups: 5
        compress: false
      trace-integration:
        enabled: false
        record-error-to-span: false
```

#### 设置隔离级别

```go
logger.SetLoggerLevel("warn")
```

#### 输出到文件

在配置文件中的 logger 选项下添加 file 项

```yaml
  logger:
    file:
      name: logs.log
      max-size: 1
      max-age: 3
      max-backups: 5
      compress: false
```

#### 自定义 logger
自定义 logger 需要实现 logger 包中的 logger 接口
```go
type Logger interface {
    Info(args ...interface{})
    Warn(args ...interface{})
    Error(args ...interface{})
    Debug(args ...interface{})
    Fatal(args ...interface{})

    Infof(fmt string, args ...interface{})
    Warnf(fmt string, args ...interface{})
    Errorf(fmt string, args ...interface{})
    Debugf(fmt string, args ...interface{})
    Fatalf(fmt string, args ...interface{})
}
```
然后调用 SetLogger 方法设置 logger
```go
logger.SetLogger(&customLogger{})
```

#### 集成 OpenTelemetry trace 信息

logger 支持 OpenTelemetry trace 集成，自动将 trace 信息注入日志。

##### 方式1：New API（推荐）

```go
ins, err := dubbo.NewInstance(
    dubbo.WithProtocol(
        protocol.WithTriple(),
        protocol.WithPort(20000),
    ),
    dubbo.WithLogger(
        log.WithZap(),
        log.WithTraceIntegration(true),
        log.WithRecordErrorToSpan(true),
    ),
)
```

##### 方式2：Old API

```go
loggerConfig := config.NewLoggerConfigBuilder().
    SetDriver("zap").
    SetTraceIntegrationEnabled(true).
    SetRecordErrorToSpan(true).
    Build()
loggerConfig.Init()
```

##### 使用 CtxLogger 记录日志

```go
import (
    "context"
)

import (
    gostLogger "github.com/dubbogo/gost/log/logger"
    "dubbo.apache.org/dubbo-go/v3/logger"
)

// 获取 CtxLogger
rawLogger := gostLogger.GetLogger()
ctxLog := rawLogger.(logger.CtxLogger)

// 创建 context（例如从请求中获取，或包含 trace 信息）
ctx := context.Background()

// 使用 context 记录日志
ctxLog.CtxInfo(ctx, "hello dubbogo this is info log")
ctxLog.CtxDebug(ctx, "hello dubbogo this is debug log")
ctxLog.CtxWarn(ctx, "hello dubbogo this is warn log")
ctxLog.CtxError(ctx, "hello dubbogo this is error log")

// 格式化日志
ctxLog.CtxInfof(ctx, "user: %s", "alice")
ctxLog.CtxDebugf(ctx, "value: %d", 42)
ctxLog.CtxWarnf(ctx, "latency: %dms", 150)
ctxLog.CtxErrorf(ctx, "failed: %v", err)
```
