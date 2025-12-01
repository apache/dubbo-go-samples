#!/bin/bash

# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

echo "Running custom generic test script"

# Check if project directory is provided
if [ -z "$1" ]; then
    echo "Project directory not provided"
    exit 1
fi

P_DIR="$1"
echo "Using project directory: $P_DIR"

# Function to clean up
cleanup() {
    echo "Cleaning up..."
    # Kill any running server processes
    if [ -f "/tmp/.generic-server.pid" ]; then
        kill $(cat /tmp/.generic-server.pid) 2>/dev/null || true
        rm -f /tmp/.generic-server.pid
    fi
}

# Set trap to ensure cleanup happens
trap cleanup EXIT INT TERM

# Build and start server
echo "Building and starting server..."
cd "$P_DIR/go-server/cmd"
go build -o server .
if [ $? -ne 0 ]; then
    echo "Failed to build server"
    exit 1
fi

# Start server in background
echo "Starting server..."
./server > server.log 2>&1 &
SERVER_PID=$!
echo $SERVER_PID > /tmp/.generic-server.pid
echo "Server started with PID: $SERVER_PID"

# Wait for server to start and register with Zookeeper
echo "Waiting for server to register with Zookeeper..."
sleep 15

# Check if server is running
if ! kill -0 $SERVER_PID 2>/dev/null; then
    echo "Server failed to start. Check server.log for details:"
    cat server.log
    exit 1
fi

# Build and run client
echo "Building and running client..."
cd "$P_DIR/go-client/cmd"
go build -o client .
if [ $? -eq 0 ]; then
    ./client
    RESULT=$?
    rm -f client
else
    echo "Failed to build client"
    RESULT=1
fi

# Print server log if client failed
if [ $RESULT -ne 0 ]; then
    echo "Client failed, printing server log..."
    echo "=== Server Log ==="
    cat "$P_DIR/go-server/cmd/server.log"
fi

# Clean up
cleanup

exit $RESULT