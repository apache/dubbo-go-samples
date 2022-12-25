
```
go install k8s.io/code-generator/cmd/go-to-protobuf@latest
go install github.com/gogo/protobuf/protoc-gen-gogo@latest
go install github.com/dubbogo/tools/cmd/protoc-gen-go-triple@v1.0.8
go install github.com/golang/protobuf/protoc-gen-go@latest

# ensure GO PATH
go mod vendor

bash rpc/triple/pb2/hack/gen-go-to-protobuf.sh

protoc \
  --proto_path=. \
  --proto_path="$GOPATH/src" \
  --go_out=rpc/triple/pb2/api \
  --go-triple_out=rpc/triple/pb2/api \
  rpc/triple/pb2/api/generated.proto
```