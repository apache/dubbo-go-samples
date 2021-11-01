export GO111MODULE="on"
export GOPROXY="http://goproxy.io"
go get -u github.com/dubbogo/tools/cmd/protoc-gen-go-triple
protoc --go_out=. --go-triple_out=. samples_api.proto