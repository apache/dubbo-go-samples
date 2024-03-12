# Dubbo-go Config_Yaml

## 1.Introduction

This example demonstrates how to configure using yaml configuration files in Dubbo-go framework

## 2.Run
```txt
.
├── go-client
│   ├── cmd
│   │   └── main.go
│   └── conf
│       └── dubbogo.yml
├── go-server
│   ├── cmd
│   │   └── main.go
│   └── conf
│       └── dubbogo.yml
└─── proto
    ├── greet.pb.go
    ├── greet.proto
    └── greet.triple.go

```
The service is defined using IDL (./proto/greet.proto) and utilizes the Triple protocol.

### build Proto
```bash
cd path_to_dubbogo-sample/config_yaml/proto
protoc --go_out=. --go-triple_out=. ./greet.proto
```
### Server
```bash
export DUBBO_GO_CONFIG_PATH="../conf/dubbogo.yml"
cd path_to_dubbogo-sample/config_yaml/go-server/cmd
go run .
```
### Client
```bash
export DUBBO_GO_CONFIG_PATH="../conf/dubbogo.yml"
cd path_to_dubbogo-sample/config_yaml/go-client/cmd
go run .
```

### 2.1 Client usage instructions

Client-defined dubbogo.yaml

```yaml
# dubbo client yaml configure file
dubbo:
  registries:
    demoZK:
      protocol: zookeeper
      timeout: 3s
      address: 127.0.0.1:2181
  consumer:
    references:
      GreetServiceImpl:
        protocol: tri
        interface: com.apache.dubbo.sample.Greeter
        registry: demoZK
        retries: 3
        timeout: 3000
```
Read and load files through `dubbo.Load()`.

```go
//...
func main() {
	//...
	if err := dubbo.Load(); err != nil {
		//...
	}
	//...
}
```

### 2.2 Server usage instructions

Server-defined dubbogo.yaml

```yaml
# dubbo server yaml configure file
dubbo:
  registries:
    demoZK:
      protocol: zookeeper
      timeout: 10s
      address: 127.0.0.1:2181
  protocols:
    tripleProtocol:
      name: tri
      port: 20000
  provider:
    services:
      GreetTripleServer:
        interface: com.apache.dubbo.sample.Greeter
```

Read and load files through `dubbo.Load()`.
```go
//...
func main() {
	//...
	if err := dubbo.Load(); err != nil {
		//...
	}
	//...
}
```
## 3.Example

### 3.1 Server 

#### IDL

Source file path ：dubbo-go-sample/context/proto/greet.proto

```protobuf
syntax = "proto3";

package greet;

option go_package = "github.com/apache/dubbo-go-samples/config_yaml/proto;greet";

message GreetRequest {
  string name = 1;
}

message GreetResponse {
  string greeting = 1;
}

service GreetService {
  rpc Greet(GreetRequest) returns (GreetResponse) {}
}
```

#### Server Handler

On the server side, define GreetTripleServer struct:
```go
type GreetServiceHandler interface {
    Greet(context.Context, *GreetRequest) (*GreetResponse, error)
}
```
Implement the GreetServiceHandler interface and register it through `greet.SetProviderService(srv common.RPCService)`

Load the configuration file through `dubbo.Load()`

Source file path ：dubbo-go-sample/config_yaml/go-server/cmd/main.go

```go

package main

import (
	"context"
	"errors"
	"fmt"

	"dubbo.apache.org/dubbo-go/v3"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	greet "github.com/apache/dubbo-go-samples/config_yaml/proto"
)

type GreetTripleServer struct {
}

func (srv *GreetTripleServer) Greet(_ context.Context, req *greet.GreetRequest) (*greet.GreetResponse, error) {
	name := req.Name
	if name != "ConfigTest" {
		errInfo := fmt.Sprintf("name is not right: %s", name)
		return nil, errors.New(errInfo)
	}

	resp := &greet.GreetResponse{Greeting: req.Name + "-Success"}
	return resp, nil
}

func main() {
	greet.SetProviderService(&GreetTripleServer{})
	if err := dubbo.Load(); err != nil {
		panic(err)
	}
	select {}
}
```

### 3.2 Client

In the client, define greet.GreetServiceImpl instance and register it through `greet.SetConsumerService(srv common.RPCService)`:

Load the configuration file through `dubbo.Load()`

Source file path ：dubbo-go-sample/config_yaml/go-client/cmd/main.go

```go
package main

import (
	"context"
	"dubbo.apache.org/dubbo-go/v3"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	greet "github.com/apache/dubbo-go-samples/config_yaml/proto"
	"github.com/dubbogo/gost/log/logger"
)

var svc = new(greet.GreetServiceImpl)

func main() {
	greet.SetConsumerService(svc)
	if err := dubbo.Load(); err != nil {
		panic(err)
	}
	req, err := svc.Greet(context.Background(), &greet.GreetRequest{Name: "ConfigTest"})
	if err != nil || req.Greeting != "ConfigTest-Success" {
		panic(err)
	}
	logger.Info("ConfigTest successfully")
}

```

### 3.3 Show

Start the server first, then the client.
If the client prints `ConfigTest Successful`, it means the configuration is loaded and the call success.
```
2024-03-11 15:47:29     INFO    cmd/main.go:39  ConfigTest successfully
```


