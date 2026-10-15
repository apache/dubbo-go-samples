#!/usr/bin/env python3
# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements. See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License. You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

"""Capture one finite run without stopping an existing registry or monitoring stack."""
import argparse
import json
import os
from pathlib import Path
import socket
import subprocess
import sys
import tempfile
import time
import urllib.request

from validate import validate

HERE = Path(__file__).resolve().parent
ROOT = HERE.parents[1]
IMAGE = "otel/opentelemetry-collector-contrib:0.104.0"


def wait_for(check, processes=(), timeout=90):
    deadline = time.monotonic() + timeout
    last_error = None
    while time.monotonic() < deadline:
        for process in processes:
            if process.poll() is not None:
                raise RuntimeError("sample exited early; inspect the captured logs")
        try:
            result = check()
            if result:
                return result
        except (OSError, ValueError, KeyError) as error:
            last_error = error
        time.sleep(0.2)
    raise RuntimeError(f"capture timed out: {last_error}; inspect the captured logs")


def fetch(port):
    with urllib.request.urlopen(f"http://127.0.0.1:{port}/prometheus", timeout=2) as response:
        return response.read()


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("output", type=Path, help="new directory for raw capture files")
    parser.add_argument("--port-base", type=int, default=24318,
                        help="four unused consecutive ports: OTLP, provider, consumer, RPC")
    args = parser.parse_args()
    output = args.output.resolve()
    output.mkdir(parents=True, exist_ok=False)
    ports = list(range(args.port_base, args.port_base + 4))
    for port in ports:
        with socket.socket() as probe:
            probe.bind(("127.0.0.1", port))
    otlp, provider_port, consumer_port, rpc_port = ports
    env = os.environ.copy()
    env.setdefault("GOMAXPROCS", "4")
    env["DUBBO_OBSERVABILITY_OTLP_ENDPOINT"] = f"127.0.0.1:{otlp}"
    env["DUBBO_OBSERVABILITY_SERVER_PORT"] = str(rpc_port)
    registry = env.get("DUBBO_OBSERVABILITY_REGISTRY_ADDRESS", "127.0.0.1:8848")
    manifest = {
        "captured_at": time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime()),
        "go": subprocess.check_output(["go", "version"], text=True).strip(),
        "sdk": json.loads(subprocess.check_output(
            ["go", "list", "-m", "-json", "dubbo.apache.org/dubbo-go/v3"], cwd=ROOT, text=True)),
        "collector_image": IMAGE,
        "registry": registry, "otlp_http": env["DUBBO_OBSERVABILITY_OTLP_ENDPOINT"],
        "metrics": {"provider": provider_port, "consumer": consumer_port},
        "rpc_port": rpc_port, "requests": 5, "successes": 4, "failures": 1,
        "group": "observability", "version": "1.0.0", "trace_integration_enabled": True,
    }
    # Avoid storing machine-local module-cache paths in a shareable capture.
    replacement = manifest["sdk"].get("Replace")
    manifest["sdk"] = {k: v for k, v in manifest["sdk"].items()
                       if k in ("Path", "Version", "Sum", "GoModSum")}
    if replacement:
        source = replacement.get("Dir")
        manifest["sdk"]["replacement"] = {key: value for key, value in replacement.items()
                                            if key in ("Version", "Sum", "GoModSum")}
        if source:
            manifest["sdk"]["replacement"]["commit"] = subprocess.check_output(
                ["git", "-C", source, "rev-parse", "HEAD"], text=True).strip()
            manifest["sdk"]["replacement"]["dirty"] = bool(subprocess.check_output(
                ["git", "-C", source, "status", "--porcelain"], text=True).strip())
    (output / "manifest.json").write_text(json.dumps(manifest, indent=2) + "\n")
    config = f"""receivers:
  otlp:
    protocols:
      http:
        endpoint: 0.0.0.0:4318
exporters:
  file:
    path: /capture/traces.jsonl
service:
  pipelines:
    traces:
      receivers: [otlp]
      exporters: [file]
"""
    (output / "collector.yaml").write_text(config)
    container = None
    processes = []
    try:
        container = subprocess.check_output([
            "docker", "run", "--detach", "--rm", "--user", "0:0",
            "--publish", f"127.0.0.1:{otlp}:4318", "--volume", f"{output}:/capture",
            IMAGE, "--config=/capture/collector.yaml"], text=True).strip()
        wait_for(lambda: socket.create_connection(("127.0.0.1", otlp), timeout=1).close() or True)
        with tempfile.TemporaryDirectory(prefix="dubbo-capture-") as build:
            binaries = {}
            for role in ("server", "client"):
                binaries[role] = str(Path(build) / role)
                subprocess.run(["go", "build", "-p", "2", "-o", binaries[role],
                                f"./observability/integration/go-{role}/cmd"],
                               cwd=ROOT, env=env, check=True)
            with (output / "provider.log").open("w") as provider_log, (output / "consumer.log").open("w") as consumer_log:
                env["DUBBO_OBSERVABILITY_METRICS_PORT"] = str(provider_port)
                provider = subprocess.Popen([binaries["server"]], cwd=ROOT, env=env,
                                            stdout=provider_log, stderr=subprocess.STDOUT)
                processes.append(provider)
                wait_for(lambda: fetch(provider_port), processes)
                # Wait for application-level registry publication before discovery.
                time.sleep(3)
                env["DUBBO_OBSERVABILITY_METRICS_PORT"] = str(consumer_port)
                consumer = subprocess.Popen([binaries["client"], "-requests", "5", "-interval", "500ms", "-hold", "10s"],
                                            cwd=ROOT, env=env, stdout=consumer_log, stderr=subprocess.STDOUT)
                processes.append(consumer)
                wait_for(lambda: "capture: requests complete" in (output / "consumer.log").read_text(), processes)
                for role, port in (("provider", provider_port), ("consumer", consumer_port)):
                    (output / f"{role}.prom").write_bytes(fetch(port))
                if consumer.wait(timeout=15) != 0:
                    raise RuntimeError("RPC contract validation failed; inspect consumer.log")
                processes.remove(consumer)
        # Every process uses the SDK exporter. Allow its final request to finish.
        result = wait_for(lambda: validate(output), processes)
        print(json.dumps(result, indent=2))
        if result["field_mismatches"]:
            raise RuntimeError("captured fields do not match; inspect field_mismatches above")
    finally:
        for process in reversed(processes):
            process.terminate()
            try:
                process.wait(timeout=5)
            except subprocess.TimeoutExpired:
                process.kill()
                process.wait()
        if container:
            with (output / "collector.log").open("w") as log:
                subprocess.run(["docker", "logs", container], stdout=log, stderr=subprocess.STDOUT)
            subprocess.run(["docker", "stop", "--time", "3", container], stdout=subprocess.DEVNULL, check=True)


if __name__ == "__main__":
    try:
        main()
    except (RuntimeError, subprocess.CalledProcessError, OSError) as error:
        print(f"capture failed: {error}", file=sys.stderr)
        sys.exit(1)
