## 日志使用

 该 samples 演示了如何使用 lumberjack 配置 dubbo-go logger

### 目录

* default:  默认打印到控制台
* level:    设置日志的隔离级别
* rolling:  输出到文件

#### 默认配置

在配置文件中不添加 logger 配置日志将会打印到控制台, 也可在配置文件中配置日志, 可参照如下方式: 

zap 日志格式和级别设置
```yaml
    logger:
    zapConfig:
      level: "debug"
      development: false
      disableCaller: false
      disableStacktrace: false
      sampling:
      encoding: "console"
    
      # encoder
      encoderConfig:
        messageKey: "message"
        levelKey: "level"
        timeKey: "time"
        nameKey: "logger"
        callerKey: "caller"
        stacktraceKey: "stacktrace"
        lineEnding: ""
        levelEncoder: "capitalColor"
        timeEncoder: "iso8601"
        durationEncoder: "seconds"
        callerEncoder: "short"
        nameEncoder: ""
    
      outputPaths:
        - "stderr"
      errorOutputPaths:
        - "stderr"
      initialFields:
```

#### 设置隔离级别

```go
logger.SetLoggerLevel("warn")
```

#### 输出到文件

在配置文件中的 logger 选项下添加 lumberjackConfig 项

```yaml
logger:
    lumberjackConfig:
      # 写日志的文件名称
      filename: "logs.log"
      # 每个日志文件长度的最大大小，单位是 MiB。默认 100MiB
      maxSize: 1
      # 日志保留的最大天数(只保留最近多少天的日志)
      maxAge: 3
      # 只保留最近多少个日志文件，用于控制程序总日志的大小
      maxBackups: 5
      # 是否使用本地时间，默认使用 UTC 时间
      localTime: true
      # 是否压缩日志文件，压缩方法 gzip
      compress: false
      # zap 配置可默认不填
    zapConfig:
```