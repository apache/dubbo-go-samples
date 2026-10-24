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


"""Reject broken correlations using captured RPC data, without a live registry.

The checked-in snapshot does not include provider.log or consumer.log.
Those files are capture.py output and match the repository *.log ignore.
This module rebuilds the lines validate() reads from traces.jsonl.
"""
import json
from pathlib import Path
import re
import shutil
import subprocess
import sys
import tempfile
import unittest

from validate import validate

FIXTURE = Path(__file__).resolve().parent / "capture"


def materialize_context_logs(capture):
    provider = []
    consumer = [
        "capture: no trace context",
        "capture: invalid trace context",
    ]
    request = 0
    for line in (capture / "traces.jsonl").read_text().splitlines():
        document = json.loads(line)
        for resource in document.get("resourceSpans", []):
            service = next((item["value"]["stringValue"]
                            for item in resource["resource"].get("attributes", [])
                            if item["key"] == "service.name"), "")
            for scope in resource.get("scopeSpans", []):
                for span in scope.get("spans", []):
                    fields = f'{{"trace_id": "{span["traceId"]}", "span_id": "{span["spanId"]}"}}'
                    if service == "dubbo-observability-server" and span["kind"] == 2 and span["name"] == "Greet":
                        provider.append(f"greet request received: name=request {fields}")
                    elif service == "dubbo-observability-client" and span["kind"] == 1:
                        request += 1
                        name = "error" if request % 5 == 0 else f"request-{request}"
                        consumer.append(f"sending greet request: name={name} {fields}")
    expected = json.loads((capture / "manifest.json").read_text())["requests"]
    if len(provider) != expected or request != expected:
        raise RuntimeError(
            f"expected {expected} context log lines per side, got provider={len(provider)} consumer={request}")
    (capture / "provider.log").write_text("\n".join(provider) + "\n")
    (capture / "consumer.log").write_text("\n".join(consumer) + "\n")


class CaptureValidationTest(unittest.TestCase):
    def setUp(self):
        self.temp = tempfile.TemporaryDirectory()
        self.addCleanup(self.temp.cleanup)
        self.capture = Path(self.temp.name) / "capture"
        shutil.copytree(FIXTURE, self.capture)
        materialize_context_logs(self.capture)

    def edit_spans(self, edit):
        path = self.capture / "traces.jsonl"
        documents = [json.loads(line) for line in path.read_text().splitlines()]
        for document in documents:
            for resource in document["resourceSpans"]:
                for scope in resource["scopeSpans"]:
                    for span in scope["spans"]:
                        edit(span)
        path.write_text("".join(json.dumps(document) + "\n" for document in documents))

    def test_capture_reports_consumer_application_mismatch(self):
        result = validate(self.capture)
        self.assertEqual(5, result["rpc_pairs"])
        self.assertEqual([{
            "side": "consumer", "field": "service.name/application_name",
            "service.name": "dubbo-observability-client",
            "application_name": "dubbo-observability-server",
        }], result["field_mismatches"])

    def test_cli_fails_for_captured_application_mismatch(self):
        run = subprocess.run([sys.executable, str(FIXTURE.parent / "validate.py"), str(self.capture)],
                             capture_output=True, text=True)
        self.assertEqual(1, run.returncode)
        self.assertEqual("consumer", json.loads(run.stdout)["field_mismatches"][0]["side"])

    def test_matching_application_has_no_mismatch(self):
        path = self.capture / "consumer.prom"
        path.write_text(path.read_text().replace(
            'application_name="dubbo-observability-server"',
            'application_name="dubbo-observability-client"'))
        self.assertEqual([], validate(self.capture)["field_mismatches"])

    def test_composite_rpc_service_fails(self):
        def edit(span):
            for attribute in span.get("attributes", []):
                if attribute["key"] == "rpc.service":
                    attribute["value"] = {"stringValue": "observability/observability.GreetService:1.0.0"}
        self.edit_spans(edit)
        mismatches = validate(self.capture)["field_mismatches"]
        self.assertEqual({"provider", "consumer"}, {
            item["side"] for item in mismatches if item["field"] == "rpc.service/interface"})
        self.assertTrue(any(item["field"] == "service.name/application_name" for item in mismatches))

    def test_missing_group_fails(self):
        path = self.capture / "provider.prom"
        path.write_text(path.read_text().replace('group="observability",', ''))
        with self.assertRaisesRegex(ValueError, "group missing"):
            validate(self.capture)

    def test_broken_parent_fails(self):
        def edit(span):
            if span["kind"] == 2:
                span["parentSpanId"] = "0000000000000000"
        self.edit_spans(edit)
        with self.assertRaisesRegex(ValueError, "SERVER parent"):
            validate(self.capture)

    def test_incorrect_failure_count_fails(self):
        path = self.capture / "consumer.prom"
        path.write_text(re.sub(r'(dubbo_consumer_requests_failed_total\{[^\n]+\}) 1', r'\g<1> 0', path.read_text()))
        with self.assertRaisesRegex(ValueError, "failure count mismatch"):
            validate(self.capture)

    def test_untraced_log_with_ids_fails(self):
        path = self.capture / "consumer.log"
        path.write_text(path.read_text().replace("capture: no trace context", 'capture: no trace context {"trace_id": "bad"}'))
        with self.assertRaisesRegex(ValueError, "unexpected trace fields"):
            validate(self.capture)

    def test_log_from_wrong_application_fails(self):
        # Keep a real exported ID, but use the consumer root in the provider log.
        consumer = next(line for line in (self.capture / "consumer.log").read_text().splitlines()
                        if "sending greet request" in line)
        ids = re.search(r'(\{"trace_id".+\})', consumer).group(1)
        path = self.capture / "provider.log"
        path.write_text(re.sub(r'\{"trace_id".+\}', ids, path.read_text(), count=1))
        with self.assertRaisesRegex(ValueError, "another application"):
            validate(self.capture)


if __name__ == "__main__":
    unittest.main()
