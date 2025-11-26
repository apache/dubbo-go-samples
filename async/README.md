# Async RPC Dubbo for Dubbo-go

[[English](README.md) | [中文](README_zh.md)]

This sample showcases how to invoke Dubbo services asynchronously with the new
`client`/`server` APIs over the Triple protocol. The client issues both a regular
callback based call (`GetUser`) and a fire-and-forget style call (`SayHello`)
while the server uses Protobuf serialization to serve Triple requests.

## Run the sample

1. **Start the provider**

   ```bash
   go run ./async/go-server/cmd/main.go
   ```

2. **Start the consumer**

   ```bash
   go run ./async/go-client/cmd/main.go
   ```

The client prints a log immediately after sending the asynchronous request and
later logs the callback result once the provider responds.
