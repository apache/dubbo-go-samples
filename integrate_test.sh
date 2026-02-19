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

set -euo pipefail

if [ -z "${1:-}" ]; then
  echo "Provide sample directory, like: ./integrate_test.sh direct"
  exit 1
fi

SAMPLE="$1"
P_DIR="$(pwd)/$SAMPLE"
PROJECT_NAME="$(basename "$P_DIR")"
GO_SERVER_LOG="/tmp/.${PROJECT_NAME}.go-server.log"
JAVA_SERVER_LOG="/tmp/.${PROJECT_NAME}.java-server.log"
PID_FILE="/tmp/.${PROJECT_NAME}.pid"
GO_CLIENT_TIMEOUT_SECONDS="${GO_CLIENT_TIMEOUT_SECONDS:-90}"
JAVA_SERVER_READY_TIMEOUT_SECONDS="${JAVA_SERVER_READY_TIMEOUT_SECONDS:-60}"
JAVA_SERVER_HOST="${JAVA_SERVER_HOST:-127.0.0.1}"
JAVA_SERVER_PORT="${JAVA_SERVER_PORT:-20000}"

if [ ! -d "$P_DIR" ]; then
  echo "Sample directory not found: $P_DIR"
  exit 1
fi

MAKE_CMD=(make PROJECT_DIR="$P_DIR" PROJECT_NAME="$PROJECT_NAME" -f Makefile)
JAVA_SERVER_RUN_SH="$(find "$P_DIR" -type f -path '*/java-server*/run.sh' -print -quit || true)"
JAVA_CLIENT_RUN_SH="$(find "$P_DIR" -type f -path '*/java-client*/run.sh' -print -quit || true)"
JAVA_SERVER_PID=""
GO_AUX_PIDS=()

JAVA_ENABLED=true
if { [ -n "$JAVA_SERVER_RUN_SH" ] || [ -n "$JAVA_CLIENT_RUN_SH" ]; } && ! command -v mvn >/dev/null 2>&1; then
  JAVA_ENABLED=false
  echo "Maven (mvn) is not available, all Java phases will be skipped for sample: $SAMPLE"
fi

run_make_target() {
  "${MAKE_CMD[@]}" "$1"
}

kill_if_running() {
  local pid="$1"
  if [ -n "$pid" ] && kill -0 "$pid" 2>/dev/null; then
    kill "$pid" 2>/dev/null || true
    sleep 1
    kill -9 "$pid" 2>/dev/null || true
  fi
}

cleanup() {
  local aux_pid
  for aux_pid in "${GO_AUX_PIDS[@]:-}"; do
    kill_if_running "$aux_pid"
  done

  kill_if_running "$JAVA_SERVER_PID"
  run_make_target stop >/dev/null 2>&1 || true
}
trap cleanup EXIT

resolve_config_path() {
  local role="$1"
  local conf_dir="$P_DIR/$role/conf"

  if [ -f "$conf_dir/dubbogo.yml" ]; then
    echo "$conf_dir/dubbogo.yml"
    return 0
  fi
  if [ -f "$conf_dir/dubbogo.yaml" ]; then
    echo "$conf_dir/dubbogo.yaml"
    return 0
  fi
  return 1
}

wait_for_tcp_port() {
  local host="$1"
  local port="$2"
  local timeout_seconds="$3"
  local elapsed=0

  while [ "$elapsed" -lt "$timeout_seconds" ]; do
    if timeout 1 bash -c "cat < /dev/null > /dev/tcp/$host/$port" >/dev/null 2>&1; then
      return 0
    fi
    sleep 1
    elapsed=$((elapsed + 1))
  done

  return 1
}

run_go_client() {
  if ! compgen -G "$P_DIR/go-client/cmd/*.go" >/dev/null; then
    echo "go-client/cmd/*.go not found in $P_DIR"
    return 1
  fi

  local client_conf
  client_conf="$(resolve_config_path "go-client" || true)"

  echo "Running Go client..."
  (
    cd "$P_DIR"
    if [ -n "$client_conf" ]; then
      export DUBBO_GO_CONFIG_PATH="$client_conf"
    fi
    go run ./go-client/cmd/*.go
  ) &
  local go_client_pid=$!
  local elapsed=0

  while kill -0 "$go_client_pid" 2>/dev/null; do
    if [ "$elapsed" -ge "$GO_CLIENT_TIMEOUT_SECONDS" ]; then
      echo "Go client timed out after ${GO_CLIENT_TIMEOUT_SECONDS}s: $SAMPLE"
      kill_if_running "$go_client_pid"
      wait "$go_client_pid" 2>/dev/null || true
      return 124
    fi
    sleep 1
    elapsed=$((elapsed + 1))
  done

  wait "$go_client_pid"
}

start_aux_go_servers() {
  local aux_server_dir
  local aux_name
  local aux_log
  local aux_pid
  local elapsed

  while IFS= read -r aux_server_dir; do
    [ -z "$aux_server_dir" ] && continue

    aux_name="$(basename "$(dirname "$aux_server_dir")")"
    aux_log="/tmp/.${PROJECT_NAME}.${aux_name}.log"

    echo "Starting auxiliary Go server: $aux_name"
    (
      cd "$P_DIR"
      go run "./${aux_server_dir#"$P_DIR"/}"/*.go
    ) >"$aux_log" 2>&1 &

    aux_pid="$!"
    GO_AUX_PIDS+=("$aux_pid")

    elapsed=0
    while kill -0 "$aux_pid" 2>/dev/null; do
      if [ "$elapsed" -ge 10 ]; then
        break
      fi
      sleep 1
      elapsed=$((elapsed + 1))
    done

    if ! kill -0 "$aux_pid" 2>/dev/null; then
      echo "Auxiliary Go server exited unexpectedly: $aux_name"
      cat "$aux_log" || true
      return 1
    fi
  done < <(find "$P_DIR" -mindepth 1 -maxdepth 1 -type d -name '*-server' ! -name 'go-server' ! -name 'java-server' -exec sh -c 'test -d "$1/cmd" && ls "$1"/cmd/*.go >/dev/null 2>&1 && echo "$1/cmd"' _ {} \;)
}

start_go_server() {
  echo "Starting Go server..."

  local server_conf
  server_conf="$(resolve_config_path "go-server" || true)"

  if [ -n "$server_conf" ]; then
    DUBBO_GO_CONFIG_PATH="$server_conf" run_make_target start >"$GO_SERVER_LOG" 2>&1
  else
    run_make_target start >"$GO_SERVER_LOG" 2>&1
  fi

  sleep 5

  if [ ! -f "$PID_FILE" ]; then
    echo "Go server pid file not found: $PID_FILE"
    cat "$GO_SERVER_LOG" || true
    return 1
  fi

  local server_pid
  server_pid="$(cat "$PID_FILE" 2>/dev/null || true)"
  if [ -z "$server_pid" ] || ! kill -0 "$server_pid" 2>/dev/null; then
    echo "Go server is not running after startup: $SAMPLE"
    cat "$GO_SERVER_LOG" || true
    local app_log="$P_DIR/go-server/dist/linux_amd64/release/${PROJECT_NAME}.log"
    [ -f "$app_log" ] && cat "$app_log" || true
    return 1
  fi
}

stop_go_server() {
  echo "Stopping Go server..."
  run_make_target stop >/dev/null 2>&1 || true
}

run_java_client_if_present() {
  if [ -z "$JAVA_CLIENT_RUN_SH" ]; then
    echo "No Java client found, skipping Java client run"
    return 0
  fi

  if [ "$JAVA_ENABLED" != "true" ]; then
    echo "Java phase disabled, skipping Java client run"
    return 0
  fi

  local java_client_dir
  java_client_dir="$(dirname "$JAVA_CLIENT_RUN_SH")"

  echo "Running Java client: $JAVA_CLIENT_RUN_SH"
  (
    cd "$java_client_dir"
    bash ./run.sh
  )
}

start_java_server_if_present() {
  if [ -z "$JAVA_SERVER_RUN_SH" ]; then
    echo "No Java server found, skipping Java server phase"
    return 1
  fi

  if [ "$JAVA_ENABLED" != "true" ]; then
    echo "Java phase disabled, skipping Java server phase"
    return 1
  fi

  local java_server_dir
  java_server_dir="$(dirname "$JAVA_SERVER_RUN_SH")"

  echo "Starting Java server: $JAVA_SERVER_RUN_SH"
  (
    cd "$java_server_dir"
    bash ./run.sh
  ) >"$JAVA_SERVER_LOG" 2>&1 &

  JAVA_SERVER_PID="$!"
  sleep 3

  if ! kill -0 "$JAVA_SERVER_PID" 2>/dev/null; then
    echo "Java server exited unexpectedly. Log:"
    cat "$JAVA_SERVER_LOG" || true
    return 1
  fi

  if ! wait_for_tcp_port "$JAVA_SERVER_HOST" "$JAVA_SERVER_PORT" "$JAVA_SERVER_READY_TIMEOUT_SECONDS"; then
    echo "Java server is running but not ready on ${JAVA_SERVER_HOST}:${JAVA_SERVER_PORT} after ${JAVA_SERVER_READY_TIMEOUT_SECONDS}s"
    cat "$JAVA_SERVER_LOG" || true
    return 1
  fi

  return 0
}

main() {
  echo "=========================================="
  echo "Starting sample flow for: $SAMPLE"
  echo "Sample directory: $P_DIR"
  echo "=========================================="

  start_go_server
  start_aux_go_servers

  run_go_client
  run_java_client_if_present

  stop_go_server

  if start_java_server_if_present; then
    run_java_client_if_present

    if ! kill -0 "$JAVA_SERVER_PID" 2>/dev/null; then
      echo "Java server exited before final Go client phase. Log:"
      cat "$JAVA_SERVER_LOG" || true
      exit 1
    fi

    run_go_client
  else
    echo "Java server phase skipped"
  fi

  echo "=========================================="
  echo "Sample flow completed for: $SAMPLE"
  echo "=========================================="
}

main
