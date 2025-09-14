# dubbo-go Helloworld 示例

本示例演示了 dubbo-go 作为 RPC 框架的基本用法。详细说明请参考[官方文档-快速开始](https://dubbo.apache.org/zh-cn/overview/mannual/golang-sdk/quickstart/)。

## 目录结构

- go-server/cmd/main.go - 服务端主程序，包含服务定义、处理器和 RPC 服务端
- go-client/cmd/main.go - RPC 客户端
- proto - API 的 protobuf 定义

## 运行方法

### 前置条件
1. 安装 `protoc` [version3][]
   参考[Protocol Buffer Compiler 安装][]。

2. 安装 `protoc-gen-go` 和 `protoc-gen-triple`
   以最新版本为例：

    ```shell
    go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31
    ```
   
    安装 protoc-gen-triple：

    ```shell
    go install github.com/dubbogo/protoc-gen-go-triple/v3@v3.0.2
    ```

3. 生成代码存根

    使用 protoc-gen-go 和 protoc-gen-go-triple 生成相关代码：

    ```shell
    protoc --go_out=. --go_opt=paths=source_relative --go-triple_out=. --go-triple_opt=paths=source_relative ./proto/greet.proto
    ```

### 启动服务端
```shell
go run ./go-server/cmd/main.go
```

测试服务端是否正常：
```shell
curl \
    --header "Content-Type: application/json" \
    --data '{"name": "Dubbo"}' \
    http://localhost:20000/greet.GreetService/Greet
```

### 启动客户端
```shell
go run ./go-client/cmd/main.go
```

[version3]: https://protobuf.dev/programming-guides/proto3/
[Protocol Buffer Compiler 安装]: https://dubbo-next.staged.apache.org/zh-cn/overview/reference/protoc-installation/

