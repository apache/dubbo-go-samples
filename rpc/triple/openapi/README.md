# Triple OpenAPI Sample

English | [中文](README_CN.md)

This sample demonstrates how to enable OpenAPI documentation for Triple protocol services in Dubbo-go. With OpenAPI enabled, the Triple server automatically generates and serves OpenAPI documentation in both JSON and YAML formats, and provides built-in Swagger UI and ReDoc pages for easy browsing and testing of your RPC services.

## What It Covers

- Enabling OpenAPI via `triple.OpenAPIEnable()` inside `triple.WithOpenAPI()`.
- Registering multiple proto-based services (including streaming) with OpenAPI support.
- Registering services with different versions and OpenAPI groups.
- Registering a non-proto (non-IDL) service alongside proto-based services.

## Contents

- `go-server/cmd/main.go` - The server that enables OpenAPI and registers multiple services.
- `go-client/cmd/main.go` - The client that calls the greet service (unary, server stream, client stream, bidi stream).
- `proto/greet/greet.proto` - Protobuf definition with unary and streaming RPCs.
- `proto/demo/demo.proto` - Protobuf definition for demonstrating versioned service registration.

## How to Run

### Start the Server

```shell
go run ./go-server/cmd/main.go
```

### Verify OpenAPI Documentation

Once the server is running, you can access the OpenAPI documentation via the following URLs:

| URL | Description |
|-----|-------------|
| `http://127.0.0.1:20000/dubbo/openapi/swagger-ui/` | Swagger UI page |
| `http://127.0.0.1:20000/dubbo/openapi/openapi.json` | OpenAPI spec in JSON format |
| `http://127.0.0.1:20000/dubbo/openapi/openapi.yaml` | OpenAPI spec in YAML format |
| `http://127.0.0.1:20000/dubbo/openapi/api-docs/default.json` | Per-group OpenAPI spec (JSON) for the `default` group |
| `http://127.0.0.1:20000/dubbo/openapi/api-docs/default.yaml` | Per-group OpenAPI spec (YAML) for the `default` group |
| `http://127.0.0.1:20000/dubbo/openapi/redoc/index.html?group=default` | ReDoc documentation for the `default` group |

For the `demo-2.0.0` group, replace `default` with `demo-2.0.0` in the URLs above.

### Run the Client

```shell
go run ./go-client/cmd/main.go
```

## Configuration Reference

### Server-side OpenAPI Configuration

`triple.WithOpenAPI()` is the entry point for OpenAPI configuration. It accepts the following options:

| Option | Description | Default |
|--------|-------------|---------|
| `triple.OpenAPIEnable()` | Enable OpenAPI documentation generation. **Required** to activate OpenAPI. | `false` |
| `triple.OpenAPIInfoTitle(title)` | Title of the OpenAPI document. | `"Dubbo-go OpenAPI"` |
| `triple.OpenAPIInfoDescription(desc)` | Description of the OpenAPI document. | `"Dubbo-go OpenAPI"` |
| `triple.OpenAPIInfoVersion(version)` | Version of the OpenAPI document. | `"1.0.0"` |
| `triple.OpenAPIPath(path)` | Base URL path for serving OpenAPI endpoints. | `"/dubbo/openapi"` |
| `triple.OpenAPIDefaultConsumesMediaTypes(types...)` | Default request content types. | `["application/json"]` |
| `triple.OpenAPIDefaultProducesMediaTypes(types...)` | Default response content types. | `["application/json"]` |
| `triple.OpenAPIDefaultHttpStatusCodes(codes...)` | Default HTTP status codes for responses. | `["200", "400", "500"]` |
| `triple.OpenAPISettings(settings)` | Additional key-value settings. | `{}` |

Example:

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

### Service-level OpenAPI Group

When registering a service, you can assign it to a specific OpenAPI group via `server.WithOpenAPIGroup()`:

```go
demo.RegisterGreetServiceHandler(srv, &DemoTripleServerV2{},
    server.WithOpenAPIGroup("demo-2.0.0"),
    server.WithVersion("2.0.0"),
)
```

Services without an explicit group fall into the `default` group.

### Services Registered in This Sample

| Service | Description |
|---------|-------------|
| `greet.GreetService` | Unary + streaming RPCs (server stream, client stream, bidi stream) |
| `demo.GreetService` (v1.0.0) | Unary RPC with version `1.0.0` |
| `demo.GreetService` (v2.0.0) | Same interface with version `2.0.0`, registered under OpenAPI group `demo-2.0.0` |
| `com.example.UserService` | Non-proto (non-IDL) service, registered without protobuf |
