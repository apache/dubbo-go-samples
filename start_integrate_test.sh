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
#array+=("game/go-server-game")
#array+=("game/go-server-gate")

# config-api
array+=("compatibility/config-api/rpc/triple")
#array+=("compatibility/config-api/configcenter/nacos")
array+=("compatibility/config-api/configcenter/zookeeper")
array+=("compatibility/config-api/config-merge")

# error
array+=("compatibility/error/triple/hessian2")
array+=("compatibility/error/triple/pb")

# metrics
array+=("compatibility/metrics")

# tracing
array+=("otel/tracing/stdout")

# direct
array+=("compatibility/direct")

# filer
array+=("filter/token")
array+=("compatibility/filter/custom")
array+=("compatibility/filter/token")

# context
array+=("compatibility/context/dubbo")
array+=("compatibility/context/triple")
array+=("context")

# registry
array+=("registry/zookeeper")
array+=("registry/nacos")

# generic
#array+=("compatibility/generic/default") # illegal service type registered

#timeout
array+=("timeout")

#healthcheck
array+=("healthcheck")

#streaming
array+=("streaming")

#retry
array+=("retry")

# rpc
array+=("compatibility/rpc/dubbo")
#array+=("compatibility/rpc/triple/codec-extension")
array+=("compatibility/rpc/triple/hessian2")
array+=("compatibility/rpc/triple/msgpack")
array+=("compatibility/rpc/triple/pb/dubbogo-grpc")
#array+=("compatibility/rpc/grpc")
array+=("compatibility/rpc/jsonrpc")
array+=("compatibility/rpc/triple/pb2")

array+=("rpc/grpc")
array+=("rpc/triple/pb")
array+=("rpc/triple/pb2")
array+=("rpc/triple/pb-json")
array+=("rpc/multi-protocols")

# tls
#array+=("compatibility/tls/dubbo")# tls.LoadX509KeyPair(certs{../../../x509/server1_cert.pem}, privateKey{../../../x509/server1_key.pem}) = err:open ../../../x509/server1_cert.pem: no such file or directory
#array+=("compatibility/tls/triple")# tls.LoadX509KeyPair(certs{../../../x509/server1_cert.pem}, privateKey{../../../x509/server1_key.pem}) = err:open ../../../x509/server1_cert.pem: no such file or directory
#array+=("compatibility/tls/grpc")# tls.LoadX509KeyPair(certs{../../../x509/server1_cert.pem}, privateKey{../../../x509/server1_key.pem}) = err:open ../../../x509/server1_cert.pem: no such file or directory

# async
array+=("compatibility/async")

# polaris
array+=("compatibility/polaris/registry")
array+=("compatibility/polaris/limit")

# error
array+=("error")

#config_center
array+=("config_center/nacos")
array+=("config_center/zookeeper")

# compatibility
## registry
array+=("compatibility/registry/zookeeper")
array+=("compatibility/registry/nacos")
array+=("compatibility/registry/etcd")
array+=("compatibility/registry/servicediscovery/zookeeper")
array+=("compatibility/registry/servicediscovery/nacos")
array+=("compatibility/registry/all/zookeeper")
array+=("compatibility/registry/all/nacos")

# config yaml
array+=("config_yaml")
# service_discovery
array+=("java_interop/service_discovery/interface")
array+=("java_interop/service_discovery/service")

# replace tls config
echo "The prefix of certificate path of the following files were replaced to \"$(pwd)/compatibility/tls\"."
find "$(pwd)/compatibility/tls" -type f -name '*.yml' -print0 | xargs -0 -n1
find "$(pwd)/compatibility/tls" -type f -name '*.yml' -print0 | xargs -0 sed -i 's#\.\.\/\.\.\/\.\.#'"$(pwd)/compatibility/tls"'#g'

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
