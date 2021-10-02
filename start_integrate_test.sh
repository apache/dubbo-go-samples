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

# helloworld
array=("helloworld")

# direct
array+=("direct")

# context
array+=("context/dubbo")

# config-api
array+=("config-api/rpc/triple")
array+=("config-api/configcenter/nacos")
array+=("config-api/configcenter/zookeeper")

# registry
array+=("registry/zookeeper")
array+=("registry/nacos")
array+=("registry/etcd")

# generic
array+=("generic/default")

# rpc
array+=("rpc/dubbo")
array+=("rpc/triple/codec-extension")
array+=("rpc/triple/hessian2")
array+=("rpc/triple/msgpack")
array+=("rpc/triple/pb/dubbogo-grpc")
array+=("rpc/grpc")
array+=("rpc/jsonrpc")

# game
array+=("game/go-server-game")
array+=("game/go-server-gate")


DOCKER_DIR=$(pwd)/integrate_test/dockercompose
docker-compose -f $DOCKER_DIR/docker-compose.yml up -d
bash -f $DOCKER_DIR/docker-health-check.sh
for((i=0;i<${#array[*]};i++))
do
	./integrate_test.sh "${array[i]}"
	result=$?
	if [ $result -gt 0 ]; then
	      docker-compose -f $DOCKER_DIR/docker-compose.yml down
        exit $result
	fi
done
docker-compose -f $DOCKER_DIR/docker-compose.yml down
