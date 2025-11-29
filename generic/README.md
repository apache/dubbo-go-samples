# Generic Call

Generic call is a mechanism that ensures information is correctly transmitted when the client does not have interface information. It generalizes POJOs into generic formats (such as dictionaries, strings), and is generally used in scenarios like integration testing and gateways.

This example demonstrates generic calls between Dubbo-Go and Dubbo Java services, showing how services can interoperate regardless of the language they're implemented in.

## Directory Structure

- go-server: Dubbo-Go server example
- go-client: Dubbo-Go client example with generic calls
- java-client: Dubbo Java client example
- java-server: Dubbo Java server example
- build: For integration test

Dubbo Java examples can be used to test interoperability with Dubbo-Go. You can start java server with go client, or go server with java client for testing.

## Prerequisites

- Docker and Docker Compose for running ZooKeeper registry
- Go 1.23+ for Dubbo-Go examples
- Java 8+ and Maven for Dubbo Java examples

## Registry

This example uses ZooKeeper as the registry. The following command starts ZooKeeper from docker, so you need to ensure that docker and docker-compose are installed first.

```shell
# Start ZooKeeper registry
docker run -d --name zookeeper -p 2181:2181 zookeeper:3.4.14
```

## Running the Examples

### Dubbo-Go Server

Using Dubbo-Go as provider, you can start it from command line tool:

```shell
cd go-server/cmd && go run server.go
```

### Dubbo-Go Client (Generic Call)

Using Dubbo-Go as consumer with generic calls:

```shell
cd go-client/cmd && go run client.go
```

### Dubbo Java Server

Using Dubbo Java as provider:

```shell
cd java-server/java-server
mvn clean package
sh run.sh
```

### Dubbo Java Client

Using Dubbo Java as consumer:

```shell
cd java-client/java-client
mvn clean package
sh run.sh
```

## Testing Interoperability

This example is designed to test interoperability between Dubbo-Go and Dubbo Java:

1. Start the ZooKeeper registry
2. Start either go-server or java-server
3. Run either go-client or java-client to test the generic calls

The client will make various generic calls to the server, including:
- GetUser1(String userId)
- GetUser2(String userId, String name)
- GetUser3(int userCode)
- GetUser4(int userCode, String name)
- GetOneUser()
- GetUsers(String[] userIdList)
- GetUsersMap(String[] userIdList)
- QueryUser(User user)
- QueryUsers(List<User> userObjectList)
- QueryAll()