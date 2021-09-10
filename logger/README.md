## Log usage

The samples demonstrate how to configure dubbo-go logger using lumberjack

### Contents

* default: print to the console by default
* level:   set the isolation level of the log
* rolling: output to file

#### print to the console by default

If you don't add a logger to the configuration file, the configuration log will be printed to the console. You can also configure the log in the configuration file. You can refer to the following method:

zap log format and level settings
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

#### set isolation level

```go
logger.SetLoggerLevel("warn")
```

#### output to file

Add the lumberjackConfig item under the logger option in the configuration file

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
