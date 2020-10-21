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

read -p "Enter integrate testing dir ( Like: direct/dubbo ) : " dir
echo "Generating auto test case for $dir "
if [ ! -d "$dir/go-client" ]; then
  mkdir "$dir/go-client"
  echo "Create $dir/go-client success. "
fi

if [ ! -d "$dir/go-server" ]; then
  mkdir "$dir/go-server"
  echo "Create $dir/go-server success. "
fi

echo "Copy travis.yml.... "
targetDir=${dir//\//\\\/}
sed "s/%DIR%/$targetDir/g" .integration/testing/.travis.yml > $dir/.travis.yml

echo "Copy integration_testing.sh.... "
cp .integration/testing/integration_testing.sh $dir/go-client/
cp .integration/testing/integration_testing.sh $dir/go-server/
echo "Auto test case for $dir finished!"