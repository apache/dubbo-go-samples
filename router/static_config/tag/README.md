# Static tag router

This example shows how to configure Dubbo-Go's tag router statically in code, without a registry or config center.

English | [中文](README_CN.md)

## Prerequisites

- Go 1.25+.

## What this sample demonstrates

- One untagged provider on port `20000`
- One `gray` tagged provider on port `20002`
- An application-scope static tag router configured by `dubbo.WithRouter(...)`

## How to run

### Run providers

Start two providers in separate terminals:

```shell
$ go run ./go-server/cmd/server.go
$ go run ./go-tag-server/cmd/server_tag.go
```

- `go-server` listens on `:20000` without a tag
- `go-tag-server` listens on `:20002` and exports the service with tag `gray`

### Run consumer

```shell
$ go run ./go-client/cmd/client.go
```

The client uses direct URLs only:

```text
tri://127.0.0.1:20000;tri://127.0.0.1:20002?dubbo.tag=gray
```

No registry or config center is required.

## Expected result

The client runs a single scenario and routes to `server-with-gray-tag`.
You should see logs similar to:

```text
INFO ... invoke successfully: receive: static tag router, response from: server-with-gray-tag
```

## Key router config

Static tag router:

```go
dubbo.WithRouter(
    router.WithScope("application"),
    router.WithKey("static-tag-provider"),
    router.WithForce(false),
    router.WithTags([]global.Tag{
        {
            Name:      "gray",
            Addresses: []string{"127.0.0.1:20002"},
        },
    }),
)
```

Request with tag attachment:

```go
ctx := context.WithValue(context.Background(), constant.AttachmentKey, map[string]string{
    constant.Tagkey: "gray",
})
```
