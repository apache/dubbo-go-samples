# Stdout 链路追踪导出器

[English](README.md) | [中文](README_zh.md)

本示例展示了 dubbo-go 使用 stdout 导出器的链路追踪功能。

## 运行方法

### 启动服务端

```shell
$ go run ./go-server/cmd/main.go
```

### 启动客户端

```shell
$ go run ./go-client/cmd/main.go
```

在服务端的控制台中，你将看到类似以下的链路追踪日志：

```shell
INFO tracing/tracing.go:53 tracing enabled, exporter: stdout
INFO tracing/tracing.go:54 tracing enabled, sampler: always_on

{
        "Name": "Greet",
        "SpanContext": {
                "TraceID": "dee1fcd3eafbcb73338aa719a9d4d4ad",
                "SpanID": "23a21f8330154882",
                "TraceFlags": "01",
                "TraceState": "",
                "Remote": false
        },
        "Parent": {
                "TraceID": "00000000000000000000000000000000",
                "SpanID": "0000000000000000",
                "TraceFlags": "00",
                "TraceState": "",
                "Remote": true
        },
        "SpanKind": 2,
        "StartTime": "2024-01-24T09:31:51.7352636+08:00",
        "EndTime": "2024-01-24T09:31:51.7352636+08:00",
        "Attributes": [
                {
                        "Key": "rpc.system",
                        "Value": {
                                "Type": "STRING",
                                "Value": "apache_dubbo"
                        }
                },
                {
                        "Key": "rpc.service",
                        "Value": {
                                "Type": "STRING",
                                "Value": "greet.GreetService"
                        }
                },
                {
                        "Key": "rpc.method",
                        "Value": {
                                "Type": "STRING",
                                "Value": "Greet"
                        }
                }
        ],
        "Events": null,
        "Links": null,
        "Status": {
                "Code": "Ok",
                "Description": ""
        },
        "DroppedAttributes": 0,
        "DroppedEvents": 0,
        "DroppedLinks": 0,
        "ChildSpanCount": 0,
        "Resource": [
                {
                        "Key": "service.name",
                        "Value": {
                                "Type": "STRING",
                                "Value": "dubbo_otel_tracing_server"
                        }
                },
                {
                        "Key": "service.namespace",
                        "Value": {
                                "Type": "STRING",
                                "Value": "dubbo-go"
                        }
                },
                {
                        "Key": "service.version",
                        "Value": {
                                "Type": "STRING",
                                "Value": ""
                        }
                }
        ],
        "InstrumentationLibrary": {
                "Name": "go.opentelemetry.io/otel",
                "Version": "v1.10.0",
                "SchemaURL": ""
        }
}

```

