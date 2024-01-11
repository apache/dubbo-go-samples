export GO111MODULE="on"
export GOPROXY="http://goproxy.io"
go install github.com/dubbogo/tools/cmd/protoc-gen-go-triple@v1.0.8
protoc --go_out=. --go-triple_out=. samples_api.proto