#!/bin/bash

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

set -e

echo "Building probe server image..."

# Build from repository root.
cd "$(dirname "$0")/../.."

# Build linux binary first. This uses local `replace => ../dubbo-go` if present.
TARGET_OS=${TARGET_OS:-linux}
TARGET_ARCH=${TARGET_ARCH:-amd64}

CGO_ENABLED=0 GOOS="$TARGET_OS" GOARCH="$TARGET_ARCH" go build -o metrics/probe/go-server/probeApp ./metrics/probe/go-server/cmd/main.go
docker build -f metrics/probe/go-server/Dockerfile -t dubbo-go-probe-server:latest .

echo "Build completed successfully."
