# Async RPC Dubbo for Dubbo-go

[[English](README.md) | [中文](README_zh.md)]

This sample showcases how to invoke Dubbo services asynchronously with the new
`client`/`server` APIs over the Triple protocol. It demonstrates both Go-to-Go
and Java-to-Go async interoperability.

## Features

- **Go Client & Server**: Async calls using `client.WithAsync()`
- **Java Client**: Async calls using `CompletableFuture` API
- **Java Server**: Async service implementation with `CompletableFuture`
- **Interoperability**: Java client can call Go server, Go client can call Java server

## Run Go to Go sample

1. **Start the Go server**

   ```bash
   go run ./async/go-server/cmd/main.go
   ```

2. **Start the Go client** (connects to Go server by default)

   ```bash
   go run ./async/go-client/cmd/main.go
   ```

The client prints "non-blocking before async callback resp: do something ... " and "test end" logs, demonstrating the non-blocking nature of async calls.

## Run Java-Go interoperability sample

This demonstrates **cross-language async calls**:

- **Go client** → **Java server**
- **Java client** → **Go server**

### Prerequisites

- Java 11 or higher
- Maven 3.6+

### Build Java modules

From the `async` directory:

```bash
mvn clean compile
```

### Test: Go client → Java server

1. **Modify the Go client URL** in `go-client/cmd/main.go`:

   ```go
   client.WithClientURL("tri://127.0.0.1:50051"),
   ```

2. **Start the Java server** (port 50051)

   ```bash
   cd java-server
   ./run.sh
   ```

3. **Start the Go client**

   ```bash
   go run ./async/go-client/cmd/main.go
   ```

The Go client will send async requests to the Java server and print "non-blocking before async callback resp: do something ... " logs.

### Test: Java client → Go server

1. **Start the Go server** (port 20000)

   ```bash
   go run ./async/go-server/cmd/main.go
   ```

2. **Start the Java client**

   ```bash
   cd java-client
   ./run.sh
   ```

The Java client will send async requests to the Go server using `CompletableFuture` callbacks.

## Port allocation

- **Go Server**: 20000
- **Java Server**: 50051

Both servers can run simultaneously without conflicts.
