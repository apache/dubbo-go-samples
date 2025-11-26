# Nacos Group Subscription Example

This example shows how to use group subscription with Nacos registry in Dubbo Go.

## How to run

1. Start Nacos server

2. Run server:
   ```bash
   cd registry/group_subscription/nacos/go-server/cmd
   go run server.go
   ```

3. Run client:
   ```bash
   cd registry/group_subscription/nacos/go-client/cmd
   go run client.go
   ```

## Key Configuration

The key configuration for group subscription is:
```go
dubbo.WithRegistry(
    registry.WithNacos(),
    registry.WithAddress("127.0.0.1:8848"),
    registry.WithGroup("test-group"), // Group subscription
),
```

Both server and client must use the same group name to discover each other.