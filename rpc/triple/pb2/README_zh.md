# dubbogo-grpc(pb2)

这个例子是 dubbo-go-grpc(pb2) 与 Triple 的一个简单例子。
它使用 `go-to-protobuf` 从 Go 结构中生成 pb2 文件。

## 内容

- api: proto files for grpc and triple respectively;
- go-server: Dubbo-go server
- go-client: Dubbo-go client
- models: models for Go server and client
- hack: hack scripts for generating

- api：分别用于grpc和triple的proto文件
- go-server：dubbo-go 服务器
- go-client：dubbo-go 的客户端
- 模型：Go 服务器和客户端的模型
- hack：用于生成的 hack 脚本

请注意，到目前为止，Triple 还不支持服务器流式 RPC 和客户端流式 RPC。


## 构建和运行

1. 安装开发工具

```shell
go install k8s.io/code-generator/cmd/go-to-protobuf@latest
go install github.com/gogo/protobuf/protoc-gen-gogo@latest
go install github.com/dubbogo/tools/cmd/protoc-gen-go-triple@v1.0.8
go install github.com/golang/protobuf/protoc-gen-go@latest
```

2. 生成 pb 文件和 Go 文件

```shell
# 注意: 确保本项目在 $GOPATH/src 里面，因为 go-to-protobuf 会使用 $GOPATH/src 作为 proto 文件的路径
# 使用 vendor 作为 proto 路径
go mod vendor

# 生成 pb 文件
bash rpc/triple/pb2/hack/gen-go-to-protobuf.sh

# 生成 RPC Go 文件 files from pb files
protoc \
  --proto_path=. \
  --proto_path="$GOPATH/src" \
  --go_out=rpc/triple/pb2/api \
  --go-triple_out=rpc/triple/pb2/api \
  rpc/triple/pb2/api/generated.proto
  
# 清理 vendor
rm -rf vendor
```

3. Run

```shell
# 启动 zk 作为注册中心
docker run --rm --name some-zookeeper -p 2181:2181 zookeeper

# 启动 server
DUBBO_GO_CONFIG_PATH=$(pwd)/rpc/triple/pb2/go-server/conf/dubbogo.yml go run rpc/triple/pb2/go-server/cmd/server.go

# 启动 client
DUBBO_GO_CONFIG_PATH=$(pwd)/rpc/triple/pb2/go-client/conf/dubbogo.yml go run rpc/triple/pb2/go-client/cmd/client.go
```