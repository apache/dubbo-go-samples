export GO111MODULE="on"
export GOPROXY="http://goproxy.io"
go get -u github.com/apache/dubbo-go/protocol/dubbo3/protoc-gen-dubbo3@3.0
protoc -I . samples_api.proto --dubbo3_out=plugins=grpc+dubbo:.