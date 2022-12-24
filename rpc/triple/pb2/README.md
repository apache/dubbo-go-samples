
```
go install k8s.io/code-generator/cmd/go-to-protobuf@latest
go install github.com/gogo/protobuf/protoc-gen-gogo@latest

# ensure GO PATH
go mod vendor

bash rpc/triple/pb2/hack/gen-go-to-protobuf.sh
```