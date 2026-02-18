# Generic Call Sample

[English](README.md) | [中文](README_zh.md)

This sample demonstrates how to use generic invocation with both Dubbo and Triple protocols for Go-Java interoperability. Generic invocation allows calling remote services without generating stubs or having the service interface locally.

## Layout

```
generic/
├── go-server/      # Go provider (Triple protocol, port 50052)
├── go-client/      # Go consumer with generic invocation (direct connection)
├── java-server/    # Java provider (Triple protocol, port 50052)
└── java-client/    # Java consumer with generic invocation
```

## Prerequisites

Start ZooKeeper (required by the server for service registration):

```bash
docker run -d --name zookeeper -p 2181:2181 zookeeper:3.8
```

## Run the Go Server

```bash
cd generic/go-server/cmd
go run .
```

The server exposes the Triple protocol on port `50052`, registers to ZooKeeper, and serves `UserProvider` with version `1.0.0` and group `triple`.

## Run the Go Client

```bash
cd generic/go-client/cmd
go run .
```

The client uses direct URL connection (`client.WithURL`) to connect to the server and performs generic calls via `cli.NewGenericService`. It tests both Dubbo protocol (port 20000) and Triple protocol (port 50052).

## Run the Java Server

Build and run from the java-server directory:

```bash
cd generic/java-server
mvn clean compile exec:java -Dexec.mainClass="org.apache.dubbo.samples.ApiProvider"
```

## Run the Java Client

```bash
cd generic/java-client
mvn clean compile exec:java -Dexec.mainClass="org.apache.dubbo.samples.ApiTripleConsumer"
```

The client uses `reference.setGeneric("true")` to perform generic calls.

## Tested Methods

| Method | Parameters | Return |
|--------|------------|--------|
| GetUser1 | String | User |
| GetUser2 | String, String | User |
| GetUser3 | int | User |
| GetUser4 | int, String | User |
| GetOneUser | - | User |
| GetUsers | String[] | User[] |
| GetUsersMap | String[] | Map<String, User> |
| QueryUser | User | User |
| QueryUsers | User[] | User[] |
| QueryAll | - | Map<String, User> |

## Expected Output

Server log:

```
Generic Go server started on port 50052
Registry: zookeeper://127.0.0.1:2181
```

Client log:

```
[Triple] GetUser1(userId string) res: {id=A003, name=Joe, age=48, ...}
[Triple] GetUser2(userId string, name string) res: {id=A003, name=lily, age=48, ...}
...
All generic call tests completed
```

## Notes

- Do NOT start Go Server and Java Server at the same time. Both listen on port 50052.
- The Go server requires ZooKeeper for service registration.
- The Go client uses direct URL connection and does not require ZooKeeper.
