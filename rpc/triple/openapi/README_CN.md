# Triple OpenAPI 示例

[English](README.md) | 中文

本示例演示如何在 Dubbo-go 中为 Triple 协议服务启用 OpenAPI 文档。启用 OpenAPI 后，Triple 服务端会自动生成并提供 JSON 和 YAML 两种格式的 OpenAPI 文档，同时内置 Swagger UI 和 ReDoc 页面，方便浏览和测试 RPC 服务。

## 功能说明

- 通过 `triple.OpenAPIEnable()` 在 `triple.WithOpenAPI()` 中启用 OpenAPI 文档。
- 注册多个基于 Protobuf 的服务（包括流式 RPC），并支持 OpenAPI 文档。
- 注册不同版本的服务，并使用不同的 OpenAPI 分组。
- 在基于 Protobuf 的服务之外，注册非 Protobuf（非 IDL）的服务。

## 目录结构

- `go-server/cmd/main.go` - 启用 OpenAPI 并注册多个服务的服务端。
- `go-client/cmd/main.go` - 调用 Greet 服务的客户端（Unary、Server Stream、Client Stream、Bidi Stream）。
- `proto/greet/greet.proto` - 包含 Unary 和流式 RPC 的 Protobuf 定义。
- `proto/demo/demo.proto` - 用于演示多版本服务注册的 Protobuf 定义。

## 运行方法

### 启动服务端

```shell
go run ./go-server/cmd/main.go
```

### 验证 OpenAPI 文档

服务启动后，可以通过以下地址访问 OpenAPI 文档：

| URL | 说明 |
|-----|------|
| `http://127.0.0.1:20000/dubbo/openapi/swagger-ui/` | Swagger UI 页面 |
| `http://127.0.0.1:20000/dubbo/openapi/openapi.json` | JSON 格式的 OpenAPI 规范 |
| `http://127.0.0.1:20000/dubbo/openapi/openapi.yaml` | YAML 格式的 OpenAPI 规范 |
| `http://127.0.0.1:20000/dubbo/openapi/api-docs/default.json` | `default` 分组的 OpenAPI 规范（JSON） |
| `http://127.0.0.1:20000/dubbo/openapi/api-docs/default.yaml` | `default` 分组的 OpenAPI 规范（YAML） |
| `http://127.0.0.1:20000/dubbo/openapi/redoc/index.html?group=default` | `default` 分组的 ReDoc 文档 |

对于 `demo-2.0.0` 分组，将上述 URL 中的 `default` 替换为 `demo-2.0.0` 即可。

### 启动客户端

```shell
go run ./go-client/cmd/main.go
```

## 配置说明

### 服务端 OpenAPI 配置

`triple.WithOpenAPI()` 是 OpenAPI 配置的入口，接受以下选项：

| 选项 | 说明 | 默认值 |
|------|------|--------|
| `triple.OpenAPIEnable()` | 启用 OpenAPI 文档生成，**必须调用**才能激活 OpenAPI。 | `false` |
| `triple.OpenAPIInfoTitle(title)` | OpenAPI 文档标题。 | `"Dubbo-go OpenAPI"` |
| `triple.OpenAPIInfoDescription(desc)` | OpenAPI 文档描述。 | `"Dubbo-go OpenAPI"` |
| `triple.OpenAPIInfoVersion(version)` | OpenAPI 文档版本。 | `"1.0.0"` |
| `triple.OpenAPIPath(path)` | OpenAPI 端点的 URL 基础路径。 | `"/dubbo/openapi"` |
| `triple.OpenAPIDefaultConsumesMediaTypes(types...)` | 默认请求内容类型。 | `["application/json"]` |
| `triple.OpenAPIDefaultProducesMediaTypes(types...)` | 默认响应内容类型。 | `["application/json"]` |
| `triple.OpenAPIDefaultHttpStatusCodes(codes...)` | 默认响应 HTTP 状态码。 | `["200", "400", "500"]` |
| `triple.OpenAPISettings(settings)` | 额外的键值对配置。 | `{}` |

示例：

```go
srv, err := server.NewServer(
    server.WithServerProtocol(
        protocol.WithTriple(
            triple.WithOpenAPI(
                triple.OpenAPIEnable(),
                triple.OpenAPIInfoTitle("OpenAPI Service"),
                triple.OpenAPIInfoDescription("A service with OpenAPI documentation"),
                triple.OpenAPIInfoVersion("1.0.0"),
            ),
        ),
        protocol.WithPort(20000),
    ),
)
```

### 服务级 OpenAPI 分组

注册服务时，可以通过 `server.WithOpenAPIGroup()` 将服务分配到指定的 OpenAPI 分组：

```go
demo.RegisterGreetServiceHandler(srv, &DemoTripleServerV2{},
    server.WithOpenAPIGroup("demo-2.0.0"),
    server.WithVersion("2.0.0"),
)
```

未显式指定分组的服务会归入 `default` 分组。

### 本示例注册的服务

| 服务 | 说明 |
|------|------|
| `greet.GreetService` | Unary + 流式 RPC（Server Stream、Client Stream、Bidi Stream） |
| `demo.GreetService` (v1.0.0) | 版本为 `1.0.0` 的 Unary RPC |
| `demo.GreetService` (v2.0.0) | 相同接口但版本为 `2.0.0`，注册在独立的 OpenAPI 分组 `demo-2.0.0` 下 |
| `com.example.UserService` | 非 Protobuf（非 IDL）服务，不依赖 Protobuf 定义 |
