.PHONY: compile
PROTOC_GEN_GO := $(GOPATH)/bin/protoc-gen-go
PROTOC := $(shell which protoc)
ifeq ($(PROTOC),)
	PROTOC = must-rebuild
endif

UNAME := $(shell uname)

$(PROTOC):
ifeq ($(UNAME), Darwin)
	brew install protobuf
endif
ifeq ($(UNAME), Linux)
	sudo apt-get install protobuf-compiler
endif

$(PROTOC_GEN_GO):
	go install github.com/dubbogo/tools/cmd/protoc-gen-dubbo3grpc@latest

helloworld.pb.go: helloworld.proto | $(PROTOC_GEN_GO) $(PROTOC)
	protoc -I . helloworld.proto --dubbo3grpc_out=plugins=grpc+dubbo3grpc:.

.PHONY: compile
compile: helloworld.pb.go

