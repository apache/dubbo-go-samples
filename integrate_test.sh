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

if [ -z "$1" ]; then
  echo "Provide test directory please, like : ./integrate_test.sh helloworld"
  exit
fi

echo "=========================================="
echo "Starting integration test for: $1"
echo "Test project directory: $(pwd)/$1"
echo "Integration test directory: $(pwd)/integrate_test/$1"
echo "=========================================="

P_DIR=$(pwd)/$1
if [ -f "$P_DIR"/build/test.sh ]; then
    echo "Found custom test script, running: $P_DIR/build/test.sh"
    "$P_DIR"/build/test.sh "$P_DIR"
    result=$?
    exit $((result))
fi

INTEGRATE_DIR=$(pwd)/integrate_test/$1

echo "Using project directory: $P_DIR"
echo "Using integration directory: $INTEGRATE_DIR"

# waiting for port release
sleep 5

# start server
echo "Starting server..."
make PROJECT_DIR=$P_DIR PROJECT_NAME=$(basename $P_DIR) INTEGRATE_DIR=$INTEGRATE_DIR -f build/Makefile start
# waiting for registry
sleep 5

# start integration
echo "Running Go integration tests..."
make PROJECT_DIR=$P_DIR PROJECT_NAME=$(basename $P_DIR) INTEGRATE_DIR=$INTEGRATE_DIR -f build/Makefile integration
result=$?

# if fail print server log
if [ $result != 0 ];then
  echo "Go integration test failed, printing server log..."
  make PROJECT_DIR=$P_DIR PROJECT_NAME=$(basename $P_DIR) INTEGRATE_DIR=$INTEGRATE_DIR -f build/Makefile print-server-log
fi

JAVA_TEST_SHELL=$INTEGRATE_DIR/tests/java
if [ -e $JAVA_TEST_SHELL ]; then
  echo "Found Java tests, running Java integration tests..."
  # run java test
  make PROJECT_DIR=$P_DIR PROJECT_NAME=$(basename $P_DIR) INTEGRATE_DIR=$INTEGRATE_DIR -f build/Makefile integration-java
  result=$?

  # if fail print server log
  if [ $result != 0 ];then
    echo "Java integration test failed, printing server log..."
    make PROJECT_DIR=$P_DIR PROJECT_NAME=$(basename $P_DIR) INTEGRATE_DIR=$INTEGRATE_DIR -f build/Makefile print-server-log
  fi
else
  echo "No Java tests found, skipping Java integration tests"
fi

echo "Cleaning up and stopping server..."
# stop server
make PROJECT_DIR=$P_DIR PROJECT_NAME=$(basename $P_DIR) INTEGRATE_DIR=$INTEGRATE_DIR -f build/Makefile clean

echo "=========================================="
echo "Integration test completed for: $1"
echo "Final result: $result"
echo "=========================================="

exit $((result))