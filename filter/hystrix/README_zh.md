# Hystrix Filter 示例

[English](README.md) | [中文](README_zh.md)

## 背景

本示例演示了如何在 dubbo-go 中使用 Hystrix filter 实现熔断器功能。Hystrix 是一个延迟和容错库，用于隔离访问远程系统、服务或第三方库的访问点，防止级联故障，实现熔断器模式。

## 实现方法

### 1. 配置 Hystrix 命令

使用 `hystrix-go` API 配置熔断器命令。资源名称格式为：
```
dubbo:consumer:InterfaceName:group:version:Method
```

**客户端配置** (`go-client/cmd/main.go`):
```go
import (
    "github.com/afex/hystrix-go/hystrix"
    _ "github.com/apache/dubbo-go-extensions/filter/hystrix"
)

func init() {
    // 资源名称格式: dubbo:consumer:接口名:分组:版本:方法
    cmdName := "dubbo:consumer:greet.GreetService:::Greet"

    hystrix.ConfigureCommand(cmdName, hystrix.CommandConfig{
        Timeout:                1000, // 超时时间（毫秒）
        MaxConcurrentRequests:  10,   // 最大并发请求数
        RequestVolumeThreshold: 5,    // 熔断器触发的最小请求数
        SleepWindow:            5000, // 熔断后尝试恢复的时间（毫秒）
        ErrorPercentThreshold:  50,   // 错误率阈值（百分比）
    })
}
```

### 2. 使用 Hystrix Filter

**客户端** (`go-client/cmd/main.go`):
```go
import (
    "dubbo.apache.org/dubbo-go/v3/client"
)

svc, err := greet.NewGreetService(cli, client.WithFilter("hystrix_consumer"))
```

## 配置参数说明

| 参数 | 说明 |
|------|------|
| `Timeout` | 命令执行超时时间（毫秒） |
| `MaxConcurrentRequests` | 同一时间最大的并发请求数 |
| `RequestVolumeThreshold` | 熔断器触发的最小请求数（滑动窗口内） |
| `SleepWindow` | 熔断器打开后，等待多久尝试恢复（毫秒） |
| `ErrorPercentThreshold` | 错误率阈值，超过此百分比熔断器打开 |

## 运行方法

### 前置条件

1. 先启动 Go 服务端，并确认 `127.0.0.1:20000` 已经可以接受请求。

### 启动服务端

```shell
cd filter/hystrix/go-server
go run ./cmd/main.go
```

### 启动客户端

```shell
cd filter/hystrix/go-client
go run ./cmd/main.go
```

## 预期输出

**客户端输出:**
```bash
=== Test 1: Sending normal requests ===
Request 1 success: Hello, request-1! (request #1)
Request 2 success: Hello, request-2! (request #2)
Request 3 success: Hello, request-3! (request #3)

=== Test 2: Sending concurrent requests ===
Concurrent request 1 success: Hello, concurrent-1! (request #4)
Concurrent request 2 success: Hello, concurrent-2! (request #5)
...

=== Test 3: Sending requests after concurrent test ===
After-test request 1 failed (circuit might be open): ...
```

当熔断器打开时，你会看到类似以下的错误：
```bash
After-test request 1 failed: hystrix: circuit open
```

## 测试熔断器功能

示例程序包含三个测试阶段：

1. **正常请求**: 发送3个正常请求，验证基本功能
2. **并发请求**: 发送15个并发请求，可能触发熔断器
3. **恢复测试**: 在并发请求后继续发送请求，观察熔断器状态

如果触发了熔断器，等待约5秒（SleepWindow配置）后再运行客户端，可以看到熔断器恢复。

## 注意事项

- Hystrix filter 主要用于**客户端**，保护调用方免受下游服务故障的影响
- 资源名称需要与实际接口、分组、版本和方法名保持一致
- 熔断器状态：关闭 → 打开 → 半开 → 关闭
- 合理配置超时时间和并发数，避免资源耗尽
