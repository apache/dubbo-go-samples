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


"""Check raw capture files. A missing field or broken correlation fails the run."""
import argparse
import json
import math
from pathlib import Path
import re
import sys


def require(ok, message):
    if not ok:
        raise ValueError(message)


def attributes(items):
    return {item["key"]: next(iter(item["value"].values())) for item in items}


def metrics(path):
    # The sample emits string labels without embedded newlines. Preserve their
    # JSON escape semantics, rather than matching label order or whole lines.
    result = {}
    for line in path.read_text().splitlines():
        match = re.fullmatch(r'(\w+)\{(.*)\} ([^ ]+)(?: \d+)?', line)
        if not match:
            continue
        name, raw, value = match.groups()
        labels = {key: json.loads(value) for key, value in
                  re.findall(r'(\w+)=("(?:[^"\\]|\\.)*")', raw)}
        if labels.get("interface") == "observability.GreetService" and labels.get("method") == "Greet":
            result.setdefault(name, []).append((labels, float(value)))
    return result


def validate(directory):
    manifest = json.loads((directory / "manifest.json").read_text())
    spans = []
    for line in (directory / "traces.jsonl").read_text().splitlines():
        for resource in json.loads(line).get("resourceSpans", []):
            service = attributes(resource["resource"].get("attributes", [])).get("service.name")
            for scope in resource.get("scopeSpans", []):
                for span in scope.get("spans", []):
                    spans.append((service, span, attributes(span.get("attributes", []))))
    require(spans, "no exported spans")
    rpc = [(service, span, attrs) for service, span, attrs in spans if span["kind"] in (2, 3) and span["name"] == "Greet"]
    expected = manifest["requests"]
    require(len(rpc) == expected * 2, f"expected {expected} CLIENT/SERVER span pairs, got {len(rpc)} spans")
    by_id = {(span["traceId"], span["spanId"]): (service, span, attrs) for service, span, attrs in spans}
    require(len(by_id) == len(spans), "duplicate span IDs in one trace")
    apps = {"provider": "dubbo-observability-server", "consumer": "dubbo-observability-client"}
    signal_summary = {}
    field_mismatches = []
    def compare(side, trace_key, trace_value, metric_key, metric_value):
        if trace_value != metric_value:
            mismatch = {"side": side, "field": f"{trace_key}/{metric_key}",
                        trace_key: trace_value, metric_key: metric_value}
            if mismatch not in field_mismatches:
                field_mismatches.append(mismatch)

    for side, app in apps.items():
        data = metrics(directory / f"{side}.prom")
        require(data, f"{side}: no Greet metrics")
        for name, samples in data.items():
            for labels, value in samples:
                compare(side, "service.name", app, "application_name", labels.get("application_name"))
                require(labels.get("group") == manifest["group"], f"{side}: group missing in {name}")
                require(labels.get("version") == manifest["version"], f"{side}: version missing in {name}")
        prefix = f"dubbo_{side}_"
        def single(suffix):
            values = data.get(prefix + suffix, [])
            require(len(values) == 1, f"{side}: expected one {suffix} series")
            value = values[0][1]
            require(math.isfinite(value), f"{side}: non-finite {suffix}")
            return value
        require(single("requests_total") == expected, f"{side}: request count mismatch")
        require(single("requests_succeed_total") == manifest["successes"], f"{side}: success count mismatch")
        require(single("requests_failed_total") == manifest["failures"], f"{side}: failure count mismatch")
        require(single("qps_total") >= 0, f"{side}: negative QPS")
        require(single("rt_milliseconds_avg") > 0, f"{side}: missing RT")
        signal_summary[side] = {"requests": expected, "successes": manifest["successes"],
                                "failures": manifest["failures"], "qps": single("qps_total"),
                                "rt_avg_ms": single("rt_milliseconds_avg")}
        role_spans = [(span, attrs) for service, span, attrs in rpc if service == app]
        require(len(role_spans) == expected, f"{side}: service.name does not match application_name")
        require(sum(span.get("status", {}).get("code") == 2 for span, _ in role_spans) == manifest["failures"], f"{side}: error span count mismatch")
        paired_clients = set()
        for span, attrs in role_spans:
            compare(side, "rpc.service", attrs.get("rpc.service"), "interface", "observability.GreetService")
            compare(side, "rpc.method", attrs.get("rpc.method"), "method", "Greet")
            require(span["kind"] == (2 if side == "provider" else 3), f"{side}: wrong span kind")
            if side == "provider":
                parent = by_id.get((span["traceId"], span.get("parentSpanId")))
                require(parent is not None and parent[0] == apps["consumer"] and parent[1]["kind"] == 3 and parent[1]["name"] == "Greet", "SERVER parent is not its CLIENT span")
                paired_clients.add((span["traceId"], span["parentSpanId"]))
        if side == "provider":
            require(len(paired_clients) == expected, "SERVER spans do not pair with distinct CLIENT spans")
        log = (directory / f"{side}.log").read_text()
        traced_lines = [line for line in log.splitlines() if
                        ("greet request received" if side == "provider" else "sending greet request") in line]
        require(len(traced_lines) == expected, f"{side}: expected {expected} context log lines")
        logged_ids = set()
        for line in traced_lines:
            trace = re.search(r'"trace_id":\s*"([0-9a-f]+)"', line)
            span_id = re.search(r'"span_id":\s*"([0-9a-f]+)"', line)
            require(trace and span_id, f"{side}: context logger IDs missing")
            logged_ids.add((trace[1], span_id[1]))
            logged_span = by_id.get((trace[1], span_id[1]))
            require(logged_span is not None, f"{side}: log IDs not in exported spans")
            require(logged_span[0] == app, f"{side}: log IDs refer to another application")
            require(logged_span[1]["kind"] == (2 if side == "provider" else 1), f"{side}: log IDs refer to the wrong span kind")
        require(len(logged_ids) == expected, f"{side}: repeated context log IDs")
    consumer = (directory / "consumer.log").read_text().splitlines()
    for marker in ("capture: no trace context", "capture: invalid trace context"):
        lines = [line for line in consumer if marker in line]
        require(len(lines) == 1, f"missing {marker} log")
        require('"trace_id"' not in lines[0] and '"span_id"' not in lines[0], f"{marker}: unexpected trace fields")
    result = {"rpc_pairs": expected, "signals": signal_summary, "context_logs": "valid, absent and invalid verified",
              "field_mismatches": field_mismatches}
    return result


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("capture", type=Path)
    args = parser.parse_args()
    try:
        result = validate(args.capture)
        print(json.dumps(result, indent=2))
        if result["field_mismatches"]:
            return 1
    except (ValueError, KeyError, OSError, json.JSONDecodeError) as error:
        print(f"validation failed: {error}", file=sys.stderr)
        return 1
    return 0


if __name__ == "__main__":
    sys.exit(main())
