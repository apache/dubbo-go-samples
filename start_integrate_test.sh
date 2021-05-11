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

# multi-registry
array+=("multi-registry/go-server")

# async
array=("async/go-server")

# attachment
array+=("attachment/go-server")

# config-api
array+=("config-api/go-server")

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

# game
#array+=("game/go-server-game")
#array+=("game/go-server-gate")

# general
array+=("general/dubbo/go-server")
array+=("general/grpc/go-server")
#array+=("general/jsonrpc/go-server")
#array+=("general/rest/go-server")

# generic
array+=("generic/go-server")

# group
array+=("group/go-server-group-a")
array+=("group/go-server-group-b")

# hello world
array+=("helloworld/go-server")

# metric
array+=("metric/go-server")

# registry
#array+=("registry/etcd/go-server")
#array+=("registry/nacos/go-server")
#array+=("registry/servicediscovery/consul/go-server")
#array+=("registry/servicediscovery/etcd/go-server")
#array+=("registry/servicediscovery/file/go-server")
#array+=("registry/servicediscovery/nacos/go-server")
array+=("registry/servicediscovery/zookeeper/go-server")

# router
#array+=("router/condition/go-server")
#array+=("router/tag/go-server")

# tls
#array+=("tls/go-server")

for((i=0;i<${#array[*]};i++))
do
	./integrate_test.sh ${array[i]}
	result=$?
	if [ $result -gt 0 ]; then
    exit $result
	fi
done

# chain
# multi-zone