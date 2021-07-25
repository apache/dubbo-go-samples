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


# Attention! when runing on Apple M1, pls start nacos&zk server on your computer first, and comment samples with # M1 ignore.

# config-api
array=("config-api/go-server")

# config-center
array+=("configcenter/apollo/go-server")
array+=("configcenter/zookeeper/go-server")
array+=("configcenter/nacos/go-server")

# context
array+=("context/go-server")

# direct
array+=("direct/go-server")

# filter
array+=("filter/custom/go-server")
array+=("filter/tpslimit/go-server")
array+=("filter/sentinel/go-server")

# game
#array+=("game/go-server-game") # Bug
#array+=("game/go-server-gate") # Bug

# general
array+=("general/dubbo/go-server")
array+=("general/grpc/go-server")
array+=("general/jsonrpc/go-server")
array+=("general/rest/go-server")

# generic
array+=("generic/go-server")

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
array+=("multi-zone/go-server-hz")
array+=("multi-zone/go-server-sh")

# registry
array+=("registry/etcd/go-server") # M1 ignore
array+=("registry/nacos/go-server")
#array+=("registry/servicediscovery/consul/go-server") # M1 ignore & Bug
array+=("registry/servicediscovery/etcd/go-server") # M1 ignore
array+=("registry/servicediscovery/file/go-server")
array+=("registry/servicediscovery/nacos/go-server")
array+=("registry/servicediscovery/zookeeper/go-server")

# router
#array+=("router/condition/go-server") # Bug
#array+=("router/tag/go-server") # Bug

# tls
array+=("tls/go-server") # Bug

# version
array+=("version/go-server-v1")
array+=("version/go-server-v2")

for((i=0;i<${#array[*]};i++))
do
	./integrate_test.sh "${array[i]}"
	result=$?
	if [ $result -gt 0 ]; then
        exit $result
	fi
done

#config-api attachement direct config-center
