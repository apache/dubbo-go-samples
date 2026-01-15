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

# helloworld
array+=("helloworld")

# game
#array+=("game/game")
#array+=("game/gate")

# tracing
array+=("otel/tracing/stdout")
array+=("otel/tracing/otlp_http_exporter")

# direct
array+=("direct")

# filter
array+=("filter/token")
array+=("filter/custom")

# registry
array+=("registry/zookeeper")
array+=("registry/nacos")
array+=("registry/polaris")

array+=("generic")

#timeout
array+=("timeout")

#healthcheck
array+=("healthcheck")

#streaming
array+=("streaming")

#retry
array+=("retry")

# rpc
array+=("rpc/grpc")
array+=("rpc/triple/pb")
array+=("rpc/triple/pb2")
array+=("rpc/triple/pb-json")
array+=("rpc/multi-protocols")

# tls
array+=("tls")

# async
array+=("async")

# error
array+=("error")

#config_center
array+=("config_center/nacos")
array+=("config_center/zookeeper")

# config yaml
array+=("config_yaml")
# service_discovery
array+=("java_interop/service_discovery/interface")
array+=("java_interop/service_discovery/service")

DOCKER_DIR=$(pwd)/integrate_test/dockercompose
DOCKER_COMPOSE_CMD="docker-compose"

if [ "$(docker compose version > /dev/null; echo $?)" -eq 0 ]; then
  DOCKER_COMPOSE_CMD="docker compose"
fi

$DOCKER_COMPOSE_CMD -f "$DOCKER_DIR"/docker-compose.yml up -d
bash -f "$DOCKER_DIR"/docker-health-check.sh
for ((i = 0; i < ${#array[*]}; i++)); do
  ./integrate_test.sh "${array[i]}"
  result=$?
  if [ $result -gt 0 ]; then
    $DOCKER_COMPOSE_CMD -f "$DOCKER_DIR"/docker-compose.yml down
    exit $result
  fi
done
$DOCKER_COMPOSE_CMD -f "$DOCKER_DIR"/docker-compose.yml down
