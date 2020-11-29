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

# async
#./integrate_test.sh async/go-server
#
## config center
#./integrate_test.sh configcenter/apollo/go-server
#./integrate_test.sh configcenter/nacos/go-server
#
## direct
#./integrate_test.sh direct/go-server
#
## filter
#./integrate_test.sh filter/custom_filter/go-server
#./integrate_test.sh filter/tpslimit/go-server

## generic
#./integrate_test.sh generic/dubbo/go-server
#
## hello world
#./integrate_test.sh helloworld/go-server
#
## metric
#./integrate_test.sh metric/go-server

array=("general/dubbo/go-server")
#通过下标遍历
for((i=0;i<${#array[*]};i++))
do
	./integrate_test.sh ${array[i]}
	if [ $? -gt 0 ]; then
    exit $?
	fi
done