GOPATH:=$(shell go env GOPATH)

.PHONY: init
init:
	@go get -u google.golang.org/protobuf/proto
	@go install github.com/golang/protobuf/protoc-gen-go@latest

.PHONY: update
update:
	@go get -u

.PHONY: tidy
tidy:
	@go mod tidy

.PHONY: build
build:
	@go build -o adservice *.go

.PHONY: docker
docker:
	@docker build -t adservice:latest .
	@docker tag adservice:latest xpunch/dubbo-go-demo-adservice:v1.0.0
	@docker push xpunch/dubbo-go-demo-adservice:v1.0.0
