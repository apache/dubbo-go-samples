# Stdout tracing exporter

This example shows dubbo-go's tracing feature with stdout exporter.

## How to run

### Run server

```shell
$ go run ./go-server/cmd/server.go
```

### Run client

```shell
$ go run ./go-client/cmd/client.go
```

In the server's console, you will see the tracing log like this:

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


