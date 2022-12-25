#!/bin/bash
# This script using protoc to generate go code from proto files.
# it support k8s.io

CURRENT_DIR=$(cd "$(dirname "$0")"; pwd)

# --apimachinery-packages string
# comma-separated list of directories to get apimachinery input types from which are needed by any API.
# Directories prefixed with '-' are not generated, directories prefixed with '+' only create types with explicit IDL instructions.
# (default "+k8s.io/apimachinery/pkg/util/intstr,+k8s.io/apimachinery/pkg/api/resource,+k8s.io/apimachinery/pkg/runtime/schema,+k8s.io/apimachinery/pkg/runtime,k8s.io/apimachinery/pkg/apis/meta/v1,k8s.io/apimachinery/pkg/apis/meta/v1beta1,k8s.io/apimachinery/pkg/apis/testapigroup/v1")
APIMACHINERY_PKGS=(
)
# temporal not supported now, because pb3 not supported now
#    go.temporal.io/api/workflow/v1

# add your go models package here
goModels=(
  github.com/apache/dubbo-go-samples/rpc/triple/pb2/models
)

packages=$(IFS=, ; echo "${goModels[*]}")

go-to-protobuf \
  --go-header-file="$CURRENT_DIR/../hack/custom-boilerplate.go.txt" \
  --packages="$packages" \
  --apimachinery-packages=$(IFS=, ; echo "${APIMACHINERY_PKGS[*]}") \
  --proto-import=./vendor
