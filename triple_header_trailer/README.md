# Dubbo-Go Triple Header and Trailer Sample

English | [中文](README_CN.md)

This sample demonstrates how to pass request metadata and read response metadata with the Triple protocol.

It covers:

- Client request headers with `triple.NewOutgoingContext` and `triple.AppendToOutgoingContext`
- Server request metadata with `triple.FromIncomingContext` and generated stream `RequestHeader`
- Server stream response headers and trailers with `ResponseHeader` / `ResponseTrailer`
- Client-side stream response headers and trailers with generated stream methods
- `http.Header` style key lookup with `Values("X-Sample-Token")`

This is different from the `context` sample. The `context` sample demonstrates Dubbo attachments through `constant.AttachmentKey`. This sample demonstrates Triple metadata APIs exposed as `http.Header`.

For generated Dubbo-Go bidi-stream and client-stream calls, pass request metadata before stream creation with `NewOutgoingContext` / `AppendToOutgoingContext`. On the server side, handlers can read the resulting metadata from `triple.FromIncomingContext`; generated stream handlers can also inspect it through `RequestHeader`.

## Run

Start the server:

```bash
go run ./triple_header_trailer/go-server/cmd
```

Run the client in another terminal:

```bash
go run ./triple_header_trailer/go-client/cmd
```

You can also run it through the integration script:

```bash
./integrate_test.sh triple_header_trailer
```

## Regenerate Code

The generated files under `proto` are produced from `proto/greet.proto`:

```bash
protoc --go_out=. --go_opt=paths=source_relative --go-triple_out=. --go-triple_opt=paths=source_relative ./triple_header_trailer/proto/greet.proto
```

## Notes

Triple metadata keys are case-insensitive. In Go, these APIs expose metadata as `http.Header`, so application code should use the normal `http.Header` accessors such as `Get` and `Values`.

The current generated unary APIs expose the business response message directly. Therefore, this sample demonstrates request and response metadata with streaming calls.
