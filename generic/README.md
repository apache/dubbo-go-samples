# Generic Sample (Triple Generic Call)

[English](README.md) | [中文](README_zh.md)

This sample demonstrates how to use generic invocation with the Triple protocol for Go-Java interoperability. Generic invocation allows calling remote services without generating stubs or having the service interface locally.

## Layout

```
generic/
├── go-server/      # Go provider listening on :50052
├── go-client/      # Go consumer with generic invocation
├── java-server/    # Java provider listening on :50052
└── java-client/    # Java consumer with generic invocation
```

## Prerequisites

Start ZooKeeper:

```bash
docker run -d --name zookeeper -p 2181:2181 zookeeper:3.8
```

## Run the Go Server

```bash
cd generic/go-server/cmd
go run .
```

The server listens on port `50052` and registers to ZooKeeper.

## Run the Go Client

```bash
cd generic/go-client/cmd
go run .
```

The client discovers the service via ZooKeeper and uses `client.WithGenericType("true")` to perform generic calls.

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
Generic Go/Java server started on port 50052
Registry: zookeeper://127.0.0.1:2181
```

Client log:

```
[PASS] GetUser1(String): {id=A003, name=Joe, age=48, ...}
[PASS] GetUser2(String, String): {id=A003, name=lily, age=48, ...}
...
[OK] All tests passed!
```

## Notes

- Do NOT start Go Server and Java Server at the same time. Both listen on port 50052.
- Make sure ZooKeeper is running before starting the server or client.
