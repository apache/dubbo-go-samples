## Log usage

The samples demonstrate how to configure dubbo-go logger using lumberjack

### Configuration

Output file settings

```yaml
lumberjackConfig:
  # The name of the log file
  filename: "logs.log"
  # The maximum size of each log file length whose unit is MiB. The default value is 100MiB.
  maxSize: 1
  # Maximum number of days to keep logs (only keep the logs of the most recent days)
  maxAge: 3
  # Only keep the most recent log files, used to control the size of the total log of the program
  maxBackups: 5
  # Whether to use local time, UTC time is used by default
  localTime: true
  # Whether to compress the log file, the compression method is gzip
  compress: false
```

Log format and level settings

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

If you do not want to output to a file, you can follow the rules below to set only zap

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

### Run

A cycle of 1000000 times is set on the go-client side, and the server will be called all the time to facilitate the viewing of the log output

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
