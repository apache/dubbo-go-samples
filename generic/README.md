# Generic Call Example

This example shows how to use generic call in Dubbo Go.

## How to run

1. Run server:
   ```bash
   cd generic/go-server/cmd
   go run server.go
   ```

2. Run client:
   ```bash
   cd generic/go-client/cmd
   go run client.go
   ```

## Key Feature

Generic call allows you to invoke remote methods without pre-generating stub code.
The interface and method calls are resolved at runtime.