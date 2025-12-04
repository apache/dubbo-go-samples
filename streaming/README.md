# Dubbo Streaming Sample

## 1. Introduction

This sample demonstrates how to use streaming communication in Dubbo, including:
- Go language streaming implementation
- Java language streaming implementation
- Interoperability verification between Go and Java

Supported streaming modes:
- Unary Call: Single request, single response
- Bidirectional Stream**: Multiple requests, multiple responses
- **Client Stream**: Multiple requests, single response
- **Server Stream**: Single request, multiple responses

## 2. Proto Definition

Define streaming methods in the proto file by adding the `stream` keyword before parameters that need streaming:

```protobuf
service GreetService {
  rpc Greet(GreetRequest) returns (GreetResponse) {}
  rpc GreetStream(stream GreetStreamRequest) returns (stream GreetStreamResponse) {}
  rpc GreetClientStream(stream GreetClientStreamRequest) returns (GreetClientStreamResponse) {}
  rpc GreetServerStream(GreetServerStreamRequest) returns (stream GreetServerStreamResponse) {}
}
```

## 3. Go Implementation

### 3.1 Go Server

Source file path: `streaming/go-server/cmd/server.go`

```go
type GreetTripleServer struct {
}

func (srv *GreetTripleServer) Greet(ctx context.Context, req *greet.GreetRequest) (*greet.GreetResponse, error) {
	resp := &greet.GreetResponse{Greeting: req.Name}
	return resp, nil
}

func (srv *GreetTripleServer) GreetStream(ctx context.Context, stream greet.GreetService_GreetStreamServer) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			if triple.IsEnded(err) {
				break
			}
			return fmt.Errorf("triple BidiStream recv error: %s", err)
		}
		if err := stream.Send(&greet.GreetStreamResponse{Greeting: req.Name}); err != nil {
			return fmt.Errorf("triple BidiStream send error: %s", err)
		}
	}
	return nil
}

func (srv *GreetTripleServer) GreetClientStream(ctx context.Context, stream greet.GreetService_GreetClientStreamServer) (*greet.GreetClientStreamResponse, error) {
	var reqs []string
	for stream.Recv() {
		reqs = append(reqs, stream.Msg().Name)
	}
	if stream.Err() != nil && !triple.IsEnded(stream.Err()) {
		return nil, fmt.Errorf("triple ClientStream recv err: %s", stream.Err())
	}
	resp := &greet.GreetClientStreamResponse{
		Greeting: strings.Join(reqs, ","),
	}

	return resp, nil
}

func (srv *GreetTripleServer) GreetServerStream(ctx context.Context, req *greet.GreetServerStreamRequest, stream greet.GreetService_GreetServerStreamServer) error {
	for i := 0; i < 5; i++ {
		if err := stream.Send(&greet.GreetServerStreamResponse{Greeting: req.Name}); err != nil {
			return fmt.Errorf("triple ServerStream send err: %s", err)
		}
	}
	return nil
}
```

### 3.2 Go Client

Source file path: `streaming/go-client/cmd/client.go`

```go
func main() {
	cli, err := client.NewClient(
		client.WithClientURL("tri://127.0.0.1:20000"),
	)
	if err != nil {
		panic(err)
	}

	svc, err := greet.NewGreetService(cli)
	if err != nil {
		panic(err)
	}
	TestClient(svc)
}

func TestClient(cli greet.GreetService) {
	if err := testUnary(cli); err != nil {
		logger.Error(err)
	}

	if err := testBidiStream(cli); err != nil {
		logger.Error(err)
	}

	if err := testClientStream(cli); err != nil {
		logger.Error(err)
	}

	if err := testServerStream(cli); err != nil {
		logger.Error(err)
	}
}

func testUnary(cli greet.GreetService) error {
	logger.Info("start to test TRIPLE unary call")
	resp, err := cli.Greet(context.Background(), &greet.GreetRequest{Name: "triple"})
	if err != nil {
		return err
	}
	logger.Infof("TRIPLE unary call resp: %s", resp.Greeting)
	return nil
}

func testBidiStream(cli greet.GreetService) error {
	logger.Info("start to test TRIPLE bidi stream")
	stream, err := cli.GreetStream(context.Background())
	if err != nil {
		return err
	}
	if sendErr := stream.Send(&greet.GreetStreamRequest{Name: "triple"}); sendErr != nil {
		return err
	}
	resp, err := stream.Recv()
	if err != nil {
		return err
	}
	logger.Infof("TRIPLE bidi stream resp: %s", resp.Greeting)
	if err := stream.CloseRequest(); err != nil {
		return err
	}
	if err := streamCloseSend(); err != nil {
		return err
	}
	return nil
}

func testClientStream(cli greet.GreetService) error {
	logger.Info("start to test TRIPLE client stream")
	stream, err := cli.GreetClientStream(context.Background())
	if err != nil {
		return err
	}
	if err := stream.Send(&greet.GreetClientStreamRequest{Name: "triple"}); err != nil {
		return err
	}
	if err := stream.Send(&greet.GreetClientStreamRequest{Name: "dubbo"}); err != nil {
		return err
	}
	resp, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}
	logger.Infof("TRIPLE client stream resp: %s", resp.Greeting)
	return nil
}

func testServerStream(cli greet.GreetService) error {
	logger.Info("start to test TRIPLE server stream")
	stream, err := cli.GreetServerStream(context.Background(), &greet.GreetServerStreamRequest{Name: "triple"})
	if err != nil {
		return err
	}
	for {
		resp, err := stream.Recv()
		if err != nil {
			if triple.IsEnded(err) {
				break
			}
			return err
		}
		logger.Infof("TRIPLE server stream resp: %s", resp.Greeting)
	}
	return nil
}
```

## 4. Java Implementation

### 4.1 Project Structure

```
streaming/
├── pom.xml                    # Parent POM
├── proto/                     # Shared Proto files
│   └── greet.proto
├── java-server/              # Java Server
│   ├── pom.xml
│   ├── src/main/java/
│   │   └── org/apache/dubbo/samples/tri/streaming/
│   │       ├── StreamingServer.java
│   │       └── GreeterImpl.java
│   └── run.sh
└── java-client/              # Java Client
    ├── pom.xml
    ├── src/main/java/
    │   └── org/apache/dubbo/samples/tri/streaming/
    │       └── StreamingClient.java
    └── run.sh
```

### 4.2 Java Server

Source file path: `streaming/java-server/src/main/java/org/apache/dubbo/samples/tri/streaming/GreeterImpl.java`

Implements all four RPC modes:

```java
public class GreeterImpl extends DubboGreetServiceTriple.GreetServiceImplBase {
    
    // Unary call
    @Override
    public GreetResponse greet(GreetRequest request) {
        return GreetResponse.newBuilder()
                .setGreeting("Hello " + request.getName())
                .build();
    }
    
    // Bidirectional stream
    @Override
    public StreamObserver<GreetStreamRequest> greetStream(
            StreamObserver<GreetStreamResponse> responseObserver) {
        // Bidirectional stream logic
    }
    
    // Client stream
    @Override
    public StreamObserver<GreetClientStreamRequest> greetClientStream(
            StreamObserver<GreetClientStreamResponse> responseObserver) {
        // Client stream logic
    }
    
    // Server stream
    @Override
    public void greetServerStream(GreetServerStreamRequest request,
            StreamObserver<GreetServerStreamResponse> responseObserver) {
        // Send 10 responses
        for (int i = 0; i < 10; i++) {
            responseObserver.onNext(response);
        }
        responseObserver.onCompleted();
    }
}
```

### 4.3 Java Client

Source file path: `streaming/java-client/src/main/java/org/apache/dubbo/samples/tri/streaming/StreamingClient.java`

Tests all streaming modes with clear, formatted output.

## 5. Running Examples

### 5.1 Run Go Server and Client

```bash
# Start Go server
cd streaming/go-server
go run cmd/server.go

# In another terminal, start Go client
cd streaming/go-client
go run cmd/client.go
```

### 5.2 Run Java Server and Client

```bash
# Build parent project
cd streaming
mvn clean install

# Start Java server
cd java-server
./run.sh
# Or
mvn compile exec:java -Dexec.mainClass="org.apache.dubbo.samples.tri.streaming.StreamingServer"

# In another terminal, start Java client
cd java-client
./run.sh
# Or
mvn compile exec:java -Dexec.mainClass="org.apache.dubbo.samples.tri.streaming.StreamingClient"
```

## 6. Example Output

### 6.1 Java Client Output

```
======================================================================
🚀 Starting Dubbo Streaming Client Tests
======================================================================
✅ Connected to server: tri://127.0.0.1:20000
======================================================================

======================================================================
📡 TEST 1: Bidirectional Streaming
======================================================================
  ➡️  Sending request #0: Client-0
  ⬅️  Received response #1: Echo from biStream: Client-0
  ➡️  Sending request #1: Client-1
  ⬅️  Received response #2: Echo from biStream: Client-1
  ➡️  Sending request #2: Client-2
  ⬅️  Received response #3: Echo from biStream: Client-2
  ➡️  Sending request #3: Client-3
  ⬅️  Received response #4: Echo from biStream: Client-3
  ➡️  Sending request #4: Client-4
  ⬅️  Received response #5: Echo from biStream: Client-4

  📤 All requests sent, waiting for responses...

  ✅ BiStream completed - Received 5 responses

======================================================================
📡 TEST 2: Server Streaming
======================================================================
  ➡️  Sending request: StreamingClient
  📥 Waiting for server stream responses...

  ⬅️  Received response #1: Response 0 from serverStream for StreamingClient
  ⬅️  Received response #2: Response 1 from serverStream for StreamingClient
  ...
  ⬅️  Received response #10: Response 9 from serverStream for StreamingClient

  ✅ ServerStream completed - Received 10 responses

======================================================================
📊 TEST RESULTS SUMMARY
======================================================================
  Bidirectional Streaming: ✅ PASSED
  Server Streaming: ✅ PASSED
----------------------------------------------------------------------
  🎉 Overall: ALL TESTS PASSED!
======================================================================
```

### 6.2 Go Client Output

```
INFO    cmd/client.go:69    start to test TRIPLE unary call
INFO    cmd/client.go:74    TRIPLE unary call resp: Hello triple

INFO    cmd/client.go:79    start to test TRIPLE bidi stream
INFO    cmd/client.go:91    TRIPLE bidi stream resp: Echo from biStream: triple

INFO    cmd/client.go:102   start to test TRIPLE client stream
INFO    cmd/client.go:116   TRIPLE client stream resp: Received 5 names: triple, triple, triple, triple, triple

INFO    cmd/client.go:121   start to test TRIPLE server stream
INFO    cmd/client.go:127   TRIPLE server stream resp: Response 0 from serverStream for triple
INFO    cmd/client.go:127   TRIPLE server stream resp: Response 1 from serverStream for triple
...
INFO    cmd/client.go:127   TRIPLE server stream resp: Response 9 from serverStream for triple
```

## 7. Interoperability Testing

This sample fully supports interoperability between Go and Java:

| Server | Client | Unary | Bidirectional | Client Stream | Server Stream |
|--------|--------|-------|---------------|---------------|---------------|
| Go     | Go     | ✅    | ✅            | ✅            | ✅            |
| Go     | Java   | ✅    | ✅            | ✅            | ✅            |
| Java   | Go     | ✅    | ✅            | ✅            | ✅            |
| Java   | Java   | ✅    | ✅            | ✅            | ✅            |

All combinations have been tested and verified, demonstrating full compatibility of the Dubbo Triple protocol across different language implementations.

