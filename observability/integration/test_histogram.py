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

"""Synthetic validator cases, not captured SDK output or end-to-end proof."""
import copy
import math
import unittest

from validate import validate_histogram

BOUNDS = [0.005, 0.01, 0.025, 0.05, 0.1, 0.3, 1]


def replace_value(data, name, value, index=0):
    labels, _ = data[name][index]
    data[name][index] = labels, value


class HistogramValidationTest(unittest.TestCase):
    def fixture(self, side="provider"):
        labels = {"application_name": f"dubbo-observability-{side}",
                  "ip": "127.0.0.1", "interface": "observability.GreetService",
                  "method": "Greet", "group": "observability", "version": "1.0.0"}
        prefix = f"dubbo_{side}_rt_milliseconds_histogram_seconds"
        # Four successes at 20 ms and one failure at 1 ms.
        data = {f"dubbo_{side}_requests_total": [(dict(labels), 5.0)],
                prefix + "_count": [(dict(labels), 5.0)],
                prefix + "_sum": [(dict(labels), 0.081)],
                prefix + "_bucket": [(dict(labels, le=str(bound)), count) for bound, count in
                                     zip(BOUNDS + [math.inf], [1., 1., 5., 5., 5., 5., 5., 5.])]}
        return prefix, data

    def test_both_sides_include_successes_and_failures(self):
        for side in ("provider", "consumer"):
            with self.subTest(side=side):
                _, data = self.fixture(side)
                validate_histogram(data, side, 5, BOUNDS)

    def test_rejects_invalid_histograms(self):
        prefix, original = self.fixture()
        cases = {
            "missing count": lambda data: data.pop(prefix + "_count"),
            "success only": lambda data: replace_value(data, prefix + "_count", 4.),
            "missing bucket": lambda data: data[prefix + "_bucket"].pop(),
            "duplicate bucket": lambda data: data[prefix + "_bucket"].append(data[prefix + "_bucket"][0]),
            "wrong bound": lambda data: data[prefix + "_bucket"][0][0].update(le="0.006"),
            "missing le": lambda data: data[prefix + "_bucket"][0][0].pop("le"),
            "wrong application": lambda data: data[prefix + "_sum"][0][0].update(application_name="other"),
            "extra label": lambda data: data[prefix + "_bucket"][0][0].update(side="provider"),
            "missing IP": lambda data: [labels.pop("ip", None) for samples in data.values() for labels, _ in samples],
            "noncumulative": lambda data: replace_value(data, prefix + "_bucket", 0., 1),
            "wrong infinity count": lambda data: replace_value(data, prefix + "_bucket", 4., -1),
            "fractional bucket": lambda data: replace_value(data, prefix + "_bucket", 0.5),
            "milliseconds sum": lambda data: replace_value(data, prefix + "_sum", 81.),
            "negative sum": lambda data: replace_value(data, prefix + "_sum", -1.),
            "nonfinite sum": lambda data: replace_value(data, prefix + "_sum", math.nan),
        }
        for name, mutate in cases.items():
            with self.subTest(name=name):
                data = copy.deepcopy(original)
                mutate(data)
                with self.assertRaises(ValueError):
                    validate_histogram(data, "provider", 5, BOUNDS)


if __name__ == "__main__":
    unittest.main()
