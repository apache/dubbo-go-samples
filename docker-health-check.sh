#!/bin/bash

set -euo pipefail

retry() {
  local name="$1"
  local attempts="$2"
  local interval="$3"
  shift 3

  local i
  for ((i = 1; i <= attempts; i++)); do
    if "$@"; then
      echo "[health-check] $name is ready"
      return 0
    fi
    echo "[health-check] waiting for $name ($i/$attempts)..."
    sleep "$interval"
  done

  echo "[health-check] $name is not ready after $attempts attempts"
  return 1
}

check_zookeeper() {
  # zookeeper 2181 is not HTTP; "empty reply" also indicates port is open.
  curl -sS --max-time 2 127.0.0.1:2181 >/dev/null 2>&1
  local code=$?
  [ "$code" -eq 0 ] || [ "$code" -eq 52 ]
}

check_nacos() {
  curl -fsS --max-time 3 http://127.0.0.1:8848/nacos/v1/console/health/liveness >/dev/null
}

check_nacos_grpc() {
  # Nacos 2.x gRPC port.
  timeout 2 bash -c 'cat < /dev/null > /dev/tcp/127.0.0.1/9848' >/dev/null 2>&1
}

check_polaris() {
  curl -fsS --max-time 3 http://127.0.0.1:8090 >/dev/null
}

check_etcd() {
  curl -fsS --max-time 3 http://127.0.0.1:2379/health >/dev/null
}

retry "zookeeper" 24 5 check_zookeeper
retry "nacos" 24 5 check_nacos
retry "nacos-grpc-9848" 24 5 check_nacos_grpc
retry "polaris" 24 5 check_polaris
retry "etcd" 24 5 check_etcd
