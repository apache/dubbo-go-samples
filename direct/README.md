# Direct Sample (Triple Direct Call)

[English](README.md) | [中文](README_zh.md)

This sample demonstrates how to use the Dubbo-Go v3 triple API to perform a point-to-point invocation without any registry. The consumer dials a target URL (`tri://127.0.0.1:20000`) directly, which makes it ideal for local debugging or traffic mirroring scenarios.

## Layout

```
direct/
├── proto/          # greet proto definition and generated triple stubs
├── go-server/      # triple provider listening on :20000
└── go-client/      # consumer dialing the provider directly
```

## Run the provider

```bash
cd direct/go-server/cmd
go run .
```

The server uses the triple protocol on port `20000` and implements `greet.GreetService`.

## Run the consumer

```bash
cd direct/go-client/cmd
go run .
```

`go-client` creates a triple client with `client.WithClientURL("tri://127.0.0.1:20000")`, so it does not require any registry or application-level configuration files.

## Expected output

Provider log:

```
INFO ... Direct server received name = dubbo-go
```

Consumer log:

```
INFO ... direct call response: hello dubbo-go
```

That's it—this sample shows the minimal code you need to stand up a direct triple connection with Dubbo-Go.

