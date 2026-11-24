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


"""Write the Dubbo dashboards embedded in dubbo-admin's Grafana ConfigMap
(read from stdin) to the directory given as the only argument."""

import json
import re
import sys

KEY = re.compile(r"^ {6}(dubbo-[a-z]+-(?:metrics|traces)\.json): \|-$")
BODY_INDENT = " " * 8
EXPECTED = 6


def main(out_dir):
    dashboards, name, body = {}, None, []
    for line in sys.stdin:
        line = line.rstrip("\r\n")
        key = KEY.match(line)
        if name and (line.startswith(BODY_INDENT) or not line.strip()):
            body.append(line[len(BODY_INDENT):])
            continue
        if name:
            dashboards[name] = "\n".join(body).strip() + "\n"
            name, body = None, []
        if key:
            name = key.group(1)
    if name:
        dashboards[name] = "\n".join(body).strip() + "\n"

    if len(dashboards) != EXPECTED:
        sys.exit(f"expected {EXPECTED} Dubbo dashboards, found {len(dashboards)}")
    for name, text in dashboards.items():
        json.loads(text)
        with open(f"{out_dir}/{name}", "w") as f:
            f.write(text)


if __name__ == "__main__":
    main(sys.argv[1])
