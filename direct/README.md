# Direct Sample (Triple Direct Call)

[English](README.md) | [中文](README_zh.md)

This sample demonstrates how to use the Dubbo-Go v3 triple API to perform a point-to-point invocation without any registry. The consumer dials a target URL (`tri://127.0.0.1:20000`) directly, which makes it ideal for local debugging or traffic mirroring scenarios.

## Layout

```
direct/
├── proto/          # greet proto definition and generated triple stubs
├── go-server/      # triple provider listening on :20000
├── go-client/      # consumer dialing the provider directly
├── java-server/    # triple provider listening on :20000
└── java-client/    # consumer dialing the provider directly
```

## Run the Golang provider

```bash
cd direct/go-server/cmd
go run .
```

The server uses the triple protocol on port `20000` and implements `greet.GreetService`.

## Run the Golang consumer

```bash
cd direct/go-client/cmd
go run .
```

`go-client` creates a triple client with `client.WithClientURL("tri://127.0.0.1:20000")`, so it does not require any registry or application-level configuration files.

## Run the Java Provider

Build all Java modules from the root directory:
(if theres no Maven in your computer,you should install [Maven][])
```shell
mvn clean compile
cd java-server
mvn exec:java "-Dexec.mainClass=org.example.server.JavaServerApp"
```

## Run the Java Consumer

Build all Java modules from the root directory:
(if theres no Maven in your computer,you should install [Maven][])
```shell
mvn clean compile
cd java-client
mvn exec:java "-Dexec.mainClass=org.example.server.JavaClientApp"
```

## Expected output

Provider log:

```
INFO ... Direct server form Golang/Java received name = Golang Client dubbo-go/Java Client dubbo
```

Consumer log:

```
INFO ... direct call response: Hello form Java/Golang Server, Golang Client dubbo-go/Java Client dubbo
```

## Attention

Do NOT Start Go Server and Java Server at the Same Time. Both the Go server and Java server listen on the same port: 20000 and expose the same Triple service path:greet.GreetService/Greet

[Maven]: https://maven.apache.org/download.cgi