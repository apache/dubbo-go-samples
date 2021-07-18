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

# tracing integrate test
array=("tracing/dubbo/go-server")

# async
array+=("async/go-server")
array+=("attachment/go-server")
array+=("config-api/go-server")
array+=("chain")
# config center
array+=("configcenter/apollo/go-server")
array+=("configcenter/nacos/go-server")
array+=("configcenter/zookeeper/go-server")

# context
array+=("context/go-server")

# direct
array+=("direct/go-server")

# filter
array+=("filter/custom/go-server")
array+=("filter/tpslimit/go-server")
array+=("filter/sentinel/go-server")

# general-dubbo
array+=("general/dubbo/go-server")

# general-dubbo3(triple)
array+=("general/dubbo3/pb/dubbogo-grpc/server/dubbogo-server")
array+=("general/dubbo3/pb/dubbogo-java/go-server")
array+=("general/dubbo3/hessian2/go-server")
array+=("general/dubbo3/msgpack/go-server")
array+=("general/dubbo3/codec-extension/go-server")

# general-grpc
array+=("general/grpc/go-server")

# generic invocation
array+=("generic/default/go-server")
#array+=("generic/protobufjson/go-server")

# group
array+=("group/go-server-group-a")
array+=("group/go-server-group-b")

# hello world
array+=("helloworld/go-server")

# metric
array+=("metric/go-server")

# multi-registry
array+=("multi-registry/go-server")

# multi-zone
array+=("multi-zone")

# registry
array+=("registry/zookeeper/go-server")
array+=("registry/etcd/go-server")
array+=("registry/nacos/go-server")

# registry/servicediscovery (app level serivce discovery)
array+=("registry/servicediscovery/zookeeper/go-server")
array+=("registry/servicediscovery/nacos/go-server")

# router integrate test can only confirm the program build success,
# the test of router logic would be fixed later
array+=("router/uniform-router/file/go-server")
array+=("router/uniform-router/file/go-server2")

for((i=0;i<${#array[*]};i++))
do
	./integrate_test.sh "${array[i]}"
	result=$?
	if [ $result -gt 0 ]; then
        exit $result
	fi
done
