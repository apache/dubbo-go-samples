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
GO_CLIENT_BIN=""
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
SAMPLE_COMPOSE_FILE=""
SAMPLE_COMPOSE_SERVICES=()
DOCKER_COMPOSE_CMD=()
OBSERVABILITY_SEMANTICS_VERIFIED=false

if [ "$SAMPLE" = "observability/integration" ]; then
  SAMPLE_COMPOSE_FILE="$P_DIR/docker-compose.yaml"
  SAMPLE_COMPOSE_SERVICES=(jaeger otel-collector prometheus grafana)
fi

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
  local server_pid=""
  local aux_pid
  for aux_pid in "${GO_AUX_PIDS[@]:-}"; do
    kill_if_running "$aux_pid"
  done

  kill_if_running "$JAVA_SERVER_PID"
  if [ -f "$PID_FILE" ]; then
    server_pid="$(cat "$PID_FILE" 2>/dev/null || true)"
    kill_if_running "$server_pid"
    rm -f "$PID_FILE"
  fi
  stop_sample_dependencies
  if [ -n "$SAMPLE_COMPOSE_FILE" ]; then
    wait_for_tcp_port_closed "127.0.0.1" "4318" 30 || true
  fi
  if [ -n "$GO_CLIENT_BIN" ]; then
    rm -f "$GO_CLIENT_BIN"
  fi
  run_make_target stop >/dev/null 2>&1 || true
}
trap cleanup EXIT

start_sample_dependencies() {
  if [ -z "$SAMPLE_COMPOSE_FILE" ]; then
    return 0
  fi

  if ! wait_for_http_url "http://127.0.0.1:8848/nacos/v1/console/health/liveness" 90 5; then
    echo "Root Nacos liveness did not remain healthy on 127.0.0.1:8848"
    return 1
  fi
  if ! wait_for_tcp_port "127.0.0.1" "9848" 60; then
    echo "Root Nacos gRPC endpoint did not become ready on 127.0.0.1:9848"
    return 1
  fi

  if docker compose version >/dev/null 2>&1; then
    DOCKER_COMPOSE_CMD=(docker compose)
  elif command -v docker-compose >/dev/null 2>&1; then
    DOCKER_COMPOSE_CMD=(docker-compose)
  else
    echo "Docker Compose is required for sample dependencies: $SAMPLE"
    return 1
  fi

  echo "Starting sample dependencies: ${SAMPLE_COMPOSE_SERVICES[*]}"
  "${DOCKER_COMPOSE_CMD[@]}" -f "$SAMPLE_COMPOSE_FILE" up -d "${SAMPLE_COMPOSE_SERVICES[@]}"

  if ! wait_for_tcp_port "127.0.0.1" "4318" 60; then
    echo "OpenTelemetry Collector did not become ready on 127.0.0.1:4318"
    return 1
  fi
  if ! wait_for_tcp_port "127.0.0.1" "9090" 60; then
    echo "Prometheus did not become ready on 127.0.0.1:9090"
    return 1
  fi
}

stop_sample_dependencies() {
  if [ -z "$SAMPLE_COMPOSE_FILE" ] || [ "${#DOCKER_COMPOSE_CMD[@]}" -eq 0 ]; then
    return 0
  fi

  "${DOCKER_COMPOSE_CMD[@]}" -f "$SAMPLE_COMPOSE_FILE" stop "${SAMPLE_COMPOSE_SERVICES[@]}" >/dev/null 2>&1 || true
  "${DOCKER_COMPOSE_CMD[@]}" -f "$SAMPLE_COMPOSE_FILE" rm -f "${SAMPLE_COMPOSE_SERVICES[@]}" >/dev/null 2>&1 || true
}

# verify_observability_semantics asserts that the telemetry pipeline actually
# received and stored the expected signals before the sample stack is torn
# down: Prometheus scrape targets are up and have scraped RPC metrics, Jaeger
# holds a cross-service trace (consumer and provider spans of the same trace),
# and Grafana reports the provisioned data source and dashboard. Any failed
# check exits non-zero so a broken scrape target, OTLP endpoint, or dashboard
# provisioning cannot silently pass the integration.
verify_observability_semantics() {
  if [ "$SAMPLE" != "observability/integration" ]; then
    return 0
  fi

  local client_pid="${1:-}"

  echo "Verifying observability telemetry semantics before teardown..."
  if ! python3 - "${client_pid:-}" <<'PY'
import base64
import json
import os
import sys
import time
import urllib.request

PROMETHEUS = "http://127.0.0.1:9090"
JAEGER = "http://127.0.0.1:16686"
GRAFANA = "http://127.0.0.1:3000"
GRAFANA_HEADERS = {
    "Authorization": "Basic " + base64.b64encode(b"admin:admin").decode("ascii"),
}


def fetch(url, headers=None):
    request = urllib.request.Request(url, headers=headers or {})
    with urllib.request.urlopen(request, timeout=2) as response:
        return response.status, response.read()


def prometheus_targets_up():
    _, body = fetch(PROMETHEUS + "/api/v1/targets")
    jobs = {}
    for target in json.loads(body)["data"].get("activeTargets", []):
        jobs.setdefault(target["labels"].get("job"), []).append(target["health"])
    for job in ("dubbo-observability-server", "dubbo-observability-client"):
        if "up" not in jobs.get(job, []):
            return False
    return True


def prometheus_metrics_scraped():
    _, body = fetch(PROMETHEUS + "/api/v1/query?query=dubbo_provider_requests_succeed_total")
    result = json.loads(body).get("data", {}).get("result", [])
    return any(float(series["value"][1]) > 0 for series in result)


def jaeger_cross_service_trace():
    _, body = fetch(JAEGER + "/api/traces?service=dubbo-observability-client&lookback=1h&limit=5")
    for trace in json.loads(body).get("data", []):
        processes = {process.get("serviceName") for process in trace.get("processes", {}).values()}
        if {"dubbo-observability-client", "dubbo-observability-server"} <= processes:
            return True
    return False


def grafana_provisioned():
    status, _ = fetch(GRAFANA + "/api/health")
    if status != 200:
        return False
    _, body = fetch(GRAFANA + "/api/datasources/uid/prometheus", GRAFANA_HEADERS)
    datasource = json.loads(body)
    if datasource.get("uid") != "prometheus" or datasource.get("type") != "prometheus":
        return False
    _, body = fetch(GRAFANA + "/api/datasources/uid/prometheus/health", GRAFANA_HEADERS)
    if json.loads(body).get("status") != "OK":
        return False
    _, body = fetch(
        GRAFANA + "/api/datasources/proxy/uid/prometheus/api/v1/query"
        "?query=dubbo_provider_requests_succeed_total",
        GRAFANA_HEADERS,
    )
    query_result = json.loads(body).get("data", {}).get("result", [])
    if not any(float(series["value"][1]) > 0 for series in query_result):
        return False
    status, _ = fetch(
        GRAFANA + "/api/dashboards/uid/dubbo-go-observability",
        GRAFANA_HEADERS,
    )
    return status == 200


checks = [
    (prometheus_targets_up, "Prometheus scrape targets are up"),
    (prometheus_metrics_scraped, "Prometheus has scraped dubbo_provider_requests_succeed_total"),
    (jaeger_cross_service_trace, "Jaeger holds a consumer/provider cross-service trace"),
    (grafana_provisioned, "Grafana data source and dashboard are provisioned"),
]
pending = checks
last_errors = {}
deadline = time.monotonic() + 90
while pending and time.monotonic() < deadline:
    if sys.argv[1]:
        try:
            os.kill(int(sys.argv[1]), 0)
        except OSError:
            print("observability client exited before semantic verification completed")
            raise SystemExit(1)
    next_pending = []
    for check, description in pending:
        try:
            if check():
                print("  ok: " + description)
                continue
        except Exception as exc:  # noqa: BLE001
            last_errors[description] = exc
        next_pending.append((check, description))
    pending = next_pending
    if pending:
        time.sleep(2)

if pending:
    for _, description in pending:
        last_error = last_errors.get(description)
        suffix = " (last error: %s)" % last_error if last_error else ""
        print("  failed: " + description + suffix)
    print("observability semantic verification failed for: " + ", ".join(
        description for _, description in pending
    ))
    raise SystemExit(1)
PY
  then
    echo "Observability telemetry semantic verification failed for: $SAMPLE"
    return 1
  fi
  echo "Observability telemetry semantics verified"
}

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
    if python3 - "$host" "$port" <<'PY' >/dev/null 2>&1
import socket
import sys

host = sys.argv[1]
port = int(sys.argv[2])

family = socket.AF_UNSPEC
type_ = socket.SOCK_STREAM

for af, socktype, proto, _, sockaddr in socket.getaddrinfo(host, port, family, type_):
    sock = None
    try:
        sock = socket.socket(af, socktype, proto)
        sock.settimeout(1.0)
        sock.connect(sockaddr)
        sys.exit(0)
    except OSError:
        continue
    finally:
        if sock is not None:
            sock.close()

sys.exit(1)
PY
    then
      return 0
    fi
    sleep 1
    elapsed=$((elapsed + 1))
  done

  return 1
}

wait_for_tcp_port_closed() {
  local host="$1"
  local port="$2"
  local timeout_seconds="$3"
  local elapsed=0

  while [ "$elapsed" -lt "$timeout_seconds" ]; do
    if python3 - "$host" "$port" <<'PY' >/dev/null 2>&1
import socket
import sys

host = sys.argv[1]
port = int(sys.argv[2])

for af, socktype, proto, _, sockaddr in socket.getaddrinfo(host, port, socket.AF_UNSPEC, socket.SOCK_STREAM):
    sock = None
    try:
        sock = socket.socket(af, socktype, proto)
        sock.settimeout(1.0)
        sock.connect(sockaddr)
        sys.exit(1)
    except OSError:
        continue
    finally:
        if sock is not None:
            sock.close()

sys.exit(0)
PY
    then
      return 0
    fi
    sleep 1
    elapsed=$((elapsed + 1))
  done

  return 1
}

wait_for_http_url() {
  local url="$1"
  local timeout_seconds="$2"
  local required_successes="${3:-1}"
  local elapsed=0
  local successes=0

  while [ "$elapsed" -lt "$timeout_seconds" ]; do
    if python3 - "$url" <<'PY' >/dev/null 2>&1
import sys
import urllib.request

try:
    with urllib.request.urlopen(sys.argv[1], timeout=1) as response:
        if response.status != 200:
            raise RuntimeError(f"unexpected HTTP status: {response.status}")
except Exception:
    sys.exit(1)
PY
    then
      successes=$((successes + 1))
      if [ "$successes" -ge "$required_successes" ]; then
        return 0
      fi
    else
      successes=0
    fi
    sleep 1
    elapsed=$((elapsed + 1))
  done

  return 1
}

wait_for_process_exit() {
  local pid="$1"
  local timeout_seconds="$2"
  local elapsed=0

  while kill -0 "$pid" 2>/dev/null; do
    if [ "$elapsed" -ge "$timeout_seconds" ]; then
      return 1
    fi
    sleep 1
    elapsed=$((elapsed + 1))
  done

  return 0
}

wait_for_log_pattern() {
  local log_file="$1"
  local pattern="$2"
  local timeout_seconds="$3"
  local elapsed=0

  while [ "$elapsed" -lt "$timeout_seconds" ]; do
    if [ -f "$log_file" ] && grep -q "$pattern" "$log_file"; then
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
  local client_args=()
  client_conf="$(resolve_config_path "go-client" || true)"

  if [ "$SAMPLE" = "observability/integration" ]; then
    # Keep the client alive while the semantic checks poll its metrics target.
    # The test stops it only after Prometheus, Jaeger, and Grafana are verified.
    client_args=(-requests 0 -interval 500ms)
    GO_CLIENT_BIN="/tmp/.${PROJECT_NAME}.go-client.bin"
  fi

  echo "Running Go client..."
  (
    cd "$P_DIR"
    if [ -n "$client_conf" ]; then
      export DUBBO_GO_CONFIG_PATH="$client_conf"
    fi
    if [ "$SAMPLE" = "observability/integration" ]; then
      go build -o "$GO_CLIENT_BIN" ./go-client/cmd
      exec "$GO_CLIENT_BIN" "${client_args[@]}"
    fi
    go run ./go-client/cmd/*.go "${client_args[@]}"
  ) &
  local go_client_pid=$!
  local elapsed=0

  if [ "$SAMPLE" = "observability/integration" ]; then
    if ! verify_observability_semantics "$go_client_pid"; then
      echo "Observability telemetry semantic verification failed for: $SAMPLE"
      kill_if_running "$go_client_pid"
      wait "$go_client_pid" 2>/dev/null || true
      return 1
    fi
    kill_if_running "$go_client_pid"
    wait "$go_client_pid" 2>/dev/null || true
    OBSERVABILITY_SEMANTICS_VERIFIED=true
    return 0
  fi

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

run_graceful_shutdown_sample() {
  local inflight_client_log="/tmp/.${PROJECT_NAME}.go-client.inflight.log"
  local reject_client_log="/tmp/.${PROJECT_NAME}.go-client.reject.log"
  local inflight_client_pid=""
  local server_pid=""
  local server_bin="/tmp/.${PROJECT_NAME}.go-server.bin"
  local client_bin="/tmp/.${PROJECT_NAME}.go-client.bin"

  if command -v lsof >/dev/null 2>&1; then
    lsof -tiTCP:20000 -sTCP:LISTEN | xargs -r kill -9 || true
  fi

  echo "Building graceful_shutdown Go server..."
  (
    cd "$P_DIR"
    go build -o "$server_bin" ./go-server/cmd
  )

  echo "Building graceful_shutdown Go client..."
  (
    cd "$P_DIR"
    go build -o "$client_bin" ./go-client/cmd
  )

  echo "Starting graceful_shutdown Go server..."
  (
    cd "$P_DIR"
    exec "$server_bin" \
      -timeout=25s \
      -step-timeout=20s \
      -consumer-update-wait=0s \
      -offline-window=0s \
      -delay=3s \
      -ignore-context-cancel=true
  ) >"$GO_SERVER_LOG" 2>&1 &
  server_pid="$!"
  echo "$server_pid" >"$PID_FILE"

  if ! wait_for_tcp_port "127.0.0.1" "20000" 30; then
    echo "graceful_shutdown server did not become ready on 127.0.0.1:20000"
    cat "$GO_SERVER_LOG" || true
    return 1
  fi

  if ! kill -0 "$server_pid" 2>/dev/null; then
    echo "graceful_shutdown server exited unexpectedly before client start"
    cat "$GO_SERVER_LOG" || true
    return 1
  fi

  echo "Running graceful_shutdown in-flight request validation..."
  (
    cd "$P_DIR"
    exec "$client_bin" \
      -addr=tri://127.0.0.1:20000 \
      -short=true \
      -max-requests=1 \
      -min-successes=1 \
      -request-timeout=10s \
      -name-prefix=integration-inflight
  ) >"$inflight_client_log" 2>&1 &
  inflight_client_pid="$!"

  if ! wait_for_log_pattern "$GO_SERVER_LOG" "Handling greet request, name=integration-inflight-1" 30; then
    echo "graceful_shutdown in-flight request did not enter the provider"
    kill_if_running "$inflight_client_pid"
    wait "$inflight_client_pid" 2>/dev/null || true
    cat "$inflight_client_log" || true
    cat "$GO_SERVER_LOG" || true
    return 1
  fi

  echo "Triggering graceful_shutdown after in-flight request entered provider..."
  kill -INT "$server_pid" 2>/dev/null || true

  if ! wait "$inflight_client_pid"; then
    echo "graceful_shutdown in-flight request validation failed"
    cat "$inflight_client_log" || true
    cat "$GO_SERVER_LOG" || true
    return 1
  fi
  inflight_client_pid=""

  if ! wait_for_log_pattern "$GO_SERVER_LOG" "Greet request finished, name=integration-inflight-1" 30; then
    echo "graceful_shutdown in-flight request did not finish in the provider"
    cat "$inflight_client_log" || true
    cat "$GO_SERVER_LOG" || true
    return 1
  fi

  if ! wait_for_process_exit "$server_pid" 30; then
    echo "graceful_shutdown server did not exit within 30s after SIGINT"
    cat "$GO_SERVER_LOG" || true
    return 1
  fi

  wait "$server_pid" 2>/dev/null || true
  server_pid=""
  rm -f "$PID_FILE"

  echo "Starting graceful_shutdown reject-stage validation server..."
  : >"$GO_SERVER_LOG"
  (
    cd "$P_DIR"
    exec "$server_bin" \
      -timeout=25s \
      -step-timeout=20s \
      -consumer-update-wait=0s \
      -offline-window=0s \
      -reject-request=true
  ) >"$GO_SERVER_LOG" 2>&1 &
  server_pid="$!"
  echo "$server_pid" >"$PID_FILE"

  if ! wait_for_tcp_port "127.0.0.1" "20000" 30; then
    echo "graceful_shutdown reject-stage server did not become ready on 127.0.0.1:20000"
    cat "$GO_SERVER_LOG" || true
    return 1
  fi

  echo "Running graceful_shutdown reject-stage probe..."
  (
    cd "$P_DIR"
    exec "$client_bin" \
      -addr=tri://127.0.0.1:20000 \
      -short=true \
      -max-requests=1 \
      -min-failures=1 \
      -request-timeout=5s \
      -name-prefix=integration-reject-probe
  ) >"$reject_client_log" 2>&1 || {
    echo "graceful_shutdown reject-stage probe failed"
    cat "$reject_client_log" || true
    cat "$GO_SERVER_LOG" || true
    return 1
  }

  if grep -q "Handling greet request, name=integration-reject-probe" "$GO_SERVER_LOG"; then
    echo "graceful_shutdown reject-stage probe reached the Greet handler"
    cat "$reject_client_log" || true
    cat "$GO_SERVER_LOG" || true
    return 1
  fi

  if ! wait_for_log_pattern "$GO_SERVER_LOG" "application is closing, new request will be rejected" 30; then
    echo "graceful_shutdown reject-stage probe was not rejected by the framework provider filter"
    cat "$reject_client_log" || true
    cat "$GO_SERVER_LOG" || true
    return 1
  fi

  kill_if_running "$server_pid"
  wait "$server_pid" 2>/dev/null || true
  server_pid=""
  rm -f "$PID_FILE"

  echo "graceful_shutdown integration completed"
}

main() {
  echo "=========================================="
  echo "Starting sample flow for: $SAMPLE"
  echo "Sample directory: $P_DIR"
  echo "=========================================="

  if [ "$SAMPLE" = "graceful_shutdown" ]; then
    run_graceful_shutdown_sample
    echo "=========================================="
    echo "Sample flow completed for: $SAMPLE"
    echo "=========================================="
    return 0
  fi

  start_go_server
  start_sample_dependencies
  start_aux_go_servers

  run_go_client
  run_java_client_if_present

  if [ "$OBSERVABILITY_SEMANTICS_VERIFIED" != "true" ]; then
    if ! verify_observability_semantics; then
      return 1
    fi
  fi

  if [ -n "$SAMPLE_COMPOSE_FILE" ]; then
    stop_sample_dependencies
    wait_for_tcp_port_closed "127.0.0.1" "4318" 30 || true
  fi
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
