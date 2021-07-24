## 日志使用

 该 samples 演示了如何使用 lumberjack 配置 dubbo-go logger

### 配置

输出文件设置

```yaml
lumberjackConfig:
  # 写日志的文件名称
  filename: "logs.log"
  # 每个日志文件长度的最大大小，默认100M
  maxSize: 1
  # 日志保留的最大天数(只保留最近多少天的日志)
  maxAge: 3
  # 只保留最近多少个日志文件，用于控制程序总日志的大小
  maxBackups: 5
  # 是否使用本地时间，默认使用UTC时间
  localTime: true
  # 是否压缩日志文件，压缩方法gzip
  compress: false
```

zap 日志格式和级别设置

```yaml
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

如果不希望输出到文件中，可以按照如下规则，仅设置 zap

```yaml
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

### 运行

在 go-client 端设置了一个 1000000 次的循环，会一直调用 server 来方便查看日志的输出情况

```go
func main() {
	config.Load()
	time.Sleep(3 * time.Second)

	for i := 0;i < 1000000; i ++ {
		test()
	}

}

func test() {
	logger.Info("\n\n\nstart to test dubbo")
	user := &pkg.User{}
	err := userProvider.GetUser(context.TODO(), []interface{}{"A001"}, user)
	if err != nil {
		logger.Infof("error: %v\n", err)
		os.Exit(1)
		return
	}
	logger.Infof("response result: %v\n", user)
}
```