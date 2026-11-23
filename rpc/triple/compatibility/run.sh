#!/usr/bin/env bash

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

set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)"
temp_dir="$(mktemp -d)"
server_pid=""
sample_dir="$repo_root"

cleanup() {
	if [[ -n "$server_pid" ]] && kill -0 "$server_pid" 2>/dev/null; then
		kill "$server_pid" 2>/dev/null || true
		wait "$server_pid" 2>/dev/null || true
	fi
	rm -rf "$temp_dir"
}
trap cleanup EXIT

for command in go grep java make mvn python3; do
	if ! command -v "$command" >/dev/null 2>&1; then
		echo "required command not found: $command" >&2
		exit 1
	fi
done

port_open() {
	python3 - <<'PYTHON' >/dev/null 2>&1
import socket

with socket.create_connection(("127.0.0.1", 20000), timeout=1):
    pass
PYTHON
}

echo "Running dubbo-go v3.3.2 client against this checkout's Triple server"
if port_open; then
	echo "TCP port 20000 is already in use" >&2
	exit 1
fi
(
	cd "$sample_dir"
	go build -mod=readonly -o "$temp_dir/current-server" ./java_interop/non-protobuf-triple/go-server/cmd/main.go
	go build -mod=readonly -o "$temp_dir/current-client" ./java_interop/non-protobuf-triple/go-client/cmd/main.go
)
(
	cd "$repo_root/rpc/triple/compatibility/legacy-client"
	GOWORK=off go build -mod=readonly -o "$temp_dir/legacy-client" .
	GOWORK=off go build -mod=readonly -o "$temp_dir/legacy-server" ./cmd/server
)
"$temp_dir/current-server" >"$temp_dir/current-server.log" 2>&1 &
server_pid=$!

ready=false
for ((attempt = 0; attempt < 60; attempt++)); do
	if ! kill -0 "$server_pid" 2>/dev/null; then
		break
	fi
	if port_open; then
		ready=true
		break
	fi
	sleep 1
done
if [[ "$ready" != true ]]; then
	cat "$temp_dir/current-server.log" >&2
	echo "current Triple server did not listen on 127.0.0.1:20000" >&2
	exit 1
fi
"$temp_dir/legacy-client"
kill "$server_pid" 2>/dev/null || true
wait "$server_pid" 2>/dev/null || true
server_pid=""

echo "Running current checkout's Go client against the dubbo-go v3.3.2 Triple server"
"$temp_dir/legacy-server" >"$temp_dir/legacy-server.log" 2>&1 &
server_pid=$!
ready=false
for ((attempt = 0; attempt < 60; attempt++)); do
	if ! kill -0 "$server_pid" 2>/dev/null; then
		break
	fi
	if port_open; then
		ready=true
		break
	fi
	sleep 1
done
if [[ "$ready" != true ]]; then
	cat "$temp_dir/legacy-server.log" >&2
	echo "legacy Triple server did not listen on 127.0.0.1:20000" >&2
	exit 1
fi
"$temp_dir/current-client" >"$temp_dir/current-client.log" 2>&1
expected_response=$'Greet response: The name of the request is:\t dubbo-go'
if ! grep -Fq -- "$expected_response" "$temp_dir/current-client.log"; then
	cat "$temp_dir/current-client.log" >&2
	echo "current client did not receive the expected response from the legacy server" >&2
	exit 1
fi
cat "$temp_dir/current-client.log"
kill "$server_pid" 2>/dev/null || true
wait "$server_pid" 2>/dev/null || true
server_pid=""

echo "Running Apache Java/Go Triple interoperability matrix"

interop_log="$temp_dir/java-go-interoperability.log"
if ! (
	cd "$sample_dir"
	bash ./integrate_test.sh java_interop/non-protobuf-triple
) >"$interop_log" 2>&1; then
	cat "$interop_log" >&2
	exit 1
fi

expected_results=(
	$'Greet response: The name of the request is:\t dubbo-go'
	"Greet response: hi, dubbo-go"
	$'Receive result ======> The name of the request is:\t world'
	"Receive result ======> hi, world"
)
for expected in "${expected_results[@]}"; do
	if ! grep -Fq -- "$expected" "$interop_log"; then
		cat "$interop_log" >&2
		echo "interoperability response not observed: $expected" >&2
		exit 1
	fi
	grep -F -- "$expected" "$interop_log"
done
