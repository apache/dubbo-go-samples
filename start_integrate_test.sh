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
array+=("game/go-server-game")
array+=("game/go-server-gate")

# config-api
array=("config-api/rpc/triple")
array+=("config-api/configcenter/nacos")
array+=("config-api/configcenter/zookeeper")
array+=("config-api/config-merge")

# error
#array+=("error/triple/hessian2") #unsupported serialization hessian2
#array+=("error/triple/pb") #error details = [type.googleapis.com/google.rpc.DebugInfo]:{stack_entries:"
                             #github.com/dubbogo/grpc-go/internal/transport.(*http2Client).reader

# metrics
array+=("metrics")

# direct
array+=("direct")

# filer
array+=("filter/custom")
array+=("filter/token")

# context
array+=("context/dubbo")
#array+=("context/triple") # ERROR   proxy_factory/default.go:146    Invoke function error: interface conversion: interface {} is nil, not []string,

# registry
#array+=("registry/zookeeper")# group and version
#array+=("registry/nacos")# group and version
#array+=("registry/etcd")# group and version
#array+=("registry/servicediscovery/zookeeper")# group and version
#array+=("registry/servicediscovery/nacos")# group and version
#array+=("registry/all/zookeeper")# group and version
#array+=("registry/all/nacos")# group and version

# generic
#array+=("generic/default") # Unsupported serialization: hessian2

# rpc
array+=("rpc/dubbo")
#array+=("rpc/triple/codec-extension") # Unsupported serialization: hessian2
#array+=("rpc/triple/hessian2") # Unsupported serialization: hessian2
#array+=("rpc/triple/msgpack") # Unsupported serialization: hessian2
#array+=("rpc/triple/pb/dubbogo-grpc") # http2: panic serving 127.0.0.1:64763: interface conversion: *triple_protocol.compatHandlerStream is not grpc.CtxSetterStream: missing method SetContext
                                        #goroutine 42 [running]:
array+=("rpc/grpc")
array+=("rpc/jsonrpc")
array+=("rpc/triple/pb2")

# tls
#array+=("tls/dubbo")# tls.LoadX509KeyPair(certs{../../../x509/server1_cert.pem}, privateKey{../../../x509/server1_key.pem}) = err:open ../../../x509/server1_cert.pem: no such file or directory
#array+=("tls/triple")# tls.LoadX509KeyPair(certs{../../../x509/server1_cert.pem}, privateKey{../../../x509/server1_key.pem}) = err:open ../../../x509/server1_cert.pem: no such file or directory
#array+=("tls/grpc")# tls.LoadX509KeyPair(certs{../../../x509/server1_cert.pem}, privateKey{../../../x509/server1_key.pem}) = err:open ../../../x509/server1_cert.pem: no such file or directory

# async
array+=("async")

# polaris
array+=("polaris/registry")
array+=("polaris/limit")

# replace tls config
echo "The prefix of certificate path of the following files were replaced to \"$(pwd)/tls\"."
find $(pwd)/tls -type f -name '*.yml' -print0 | xargs -0 -n1
find $(pwd)/tls -type f -name '*.yml' -print0 | xargs -0 sed -i  's#\.\.\/\.\.\/\.\.#'"$(pwd)/tls"'#g'

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
