#!/bin/bash

#
#  Licensed to the Apache Software Foundation (ASF) under one or more
#  contributor license agreements.  See the NOTICE file distributed with
#  this work for additional information regarding copyright ownership.
#  The ASF licenses this file to You under the Apache License, Version 2.0
#  (the "License"); you may not use this file except in compliance with
#  the License.  You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
#  Unless required by applicable law or agreed to in writing, software
#  distributed under the License is distributed on an "AS IS" BASIS,
#  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#  See the License for the specific language governing permissions and
#  limitations under the License.
#

set -euo pipefail

# helloworld
array+=("helloworld")

# tracing
array+=("otel/tracing/stdout")

# direct
array+=("direct")

# filter
array+=("filter/token")
array+=("filter/custom")

# registry
array+=("registry/zookeeper")
array+=("registry/nacos")
array+=("registry/etcd")
array+=("registry/polaris")

array+=("generic")

# timeout
array+=("timeout")

# healthcheck
array+=("healthcheck")

# streaming
array+=("streaming")

# retry
array+=("retry")

# rpc
array+=("rpc/grpc")
array+=("rpc/multi-protocols")
array+=("rpc/triple/instance")
array+=("rpc/triple/old_triple")
array+=("rpc/triple/pb")
array+=("rpc/triple/pb-json")
array+=("rpc/triple/pb2")
array+=("rpc/triple/reflection")
array+=("rpc/triple/registry")
array+=("rpc/triple/stream")

# tls
array+=("tls")

# async
array+=("async")

# error
array+=("error")

# config_center
array+=("config_center/nacos")
# array+=("config_center/apollo")
array+=("config_center/zookeeper")

# config yaml
array+=("config_yaml")

# service_discovery
array+=("java_interop/non-protobuf-dubbo")
array+=("java_interop/non-protobuf-triple")
array+=("java_interop/protobuf-triple")
array+=("java_interop/service_discovery/interface")
array+=("java_interop/service_discovery/service")

DOCKER_DIR="$(pwd)"
DOCKER_COMPOSE_CMD="docker-compose"

if [ "$(docker compose version > /dev/null; echo $?)" -eq 0 ]; then
  DOCKER_COMPOSE_CMD="docker compose"
fi

cleanup() {
  echo "::group::> docker down"
  $DOCKER_COMPOSE_CMD -f "$DOCKER_DIR"/docker-compose.yml down >/dev/null 2>&1 || true
  echo "::endgroup::"
}
trap cleanup EXIT

echo "::group::> docker up"
$DOCKER_COMPOSE_CMD -f "$DOCKER_DIR"/docker-compose.yml up -d
echo "::endgroup::"

echo "::group::> docker health-check"
bash -f "$DOCKER_DIR"/docker-health-check.sh
echo "::endgroup::"

for t in "${array[@]}"; do
  echo "::group::> start: $t"
  set +e
  bash ./integrate_test.sh "$t"
  result=$?
  set -e
  echo "::endgroup::"

  if [ "$result" -ne 0 ]; then
    echo "[ERROR] failed: $t (exit code: $result)"
    exit "$result"
  fi

  echo "> ok: $t"
done
