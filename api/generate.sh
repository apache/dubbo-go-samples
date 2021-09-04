export GO111MODULE="on"
export GOPROXY="http://goproxy.io"
go get -u github.com/dubbogo/tools/cmd/protoc-gen-triple
protoc -I . samples_api.proto --triple_out=plugins=triple:.