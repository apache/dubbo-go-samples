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

"""SDK provenance must work for both remote modules and local source trees."""
import copy
import subprocess
import unittest
from unittest.mock import patch

from capture import sdk_metadata


class SDKMetadataTest(unittest.TestCase):
    def module(self):
        return {"Path": "dubbo.apache.org/dubbo-go/v3", "Version": "v3.3.2",
                "Sum": "module-sum", "GoModSum": "mod-sum",
                "Dir": "/private/cache/sdk", "GoMod": "/private/cache/sdk/go.mod"}

    def test_unreplaced_module_omits_machine_paths_without_mutation(self):
        module = self.module()
        original = copy.deepcopy(module)
        with patch("capture.subprocess.run") as run:
            result = sdk_metadata(module)
        run.assert_not_called()
        self.assertEqual({key: original[key] for key in ("Path", "Version", "Sum", "GoModSum")}, result)
        self.assertEqual(original, module)

    def test_versioned_replacement_never_runs_git_in_module_cache(self):
        module = self.module()
        module["Replace"] = {"Path": "github.com/example/dubbo-go/v3",
                             "Version": "v3.0.0-20261112014637-221c96e2271e",
                             "Sum": "fork-sum", "GoModSum": "fork-mod-sum",
                             "Dir": "/private/cache/fork", "GoMod": "/private/cache/fork/go.mod"}
        with patch("capture.subprocess.run") as run:
            result = sdk_metadata(module)
        run.assert_not_called()
        self.assertEqual({key: module["Replace"][key] for key in ("Path", "Version", "Sum", "GoModSum")},
                         result["replacement"])

    def test_local_git_replacement_records_revision_and_dirty_state(self):
        for dirty in (False, True):
            with self.subTest(dirty=dirty):
                module = self.module()
                module["Replace"] = {"Path": "../sdk", "Dir": "/private/sdk"}
                replies = [subprocess.CompletedProcess([], 0, "abc123\n", ""),
                           subprocess.CompletedProcess([], 0, " M metrics/options.go\n" if dirty else "", "")]
                with patch("capture.subprocess.run", side_effect=replies) as run:
                    result = sdk_metadata(module)
                self.assertEqual({"local": True, "commit": "abc123", "dirty": dirty}, result["replacement"])
                self.assertEqual(["git", "-C", "/private/sdk", "rev-parse", "HEAD"], run.call_args_list[0].args[0])

    def test_local_non_git_source_is_not_reported_as_clean_checkout(self):
        module = self.module()
        module["Replace"] = {"Path": "../sdk", "Dir": "/private/sdk"}
        with patch("capture.subprocess.run", return_value=subprocess.CompletedProcess([], 128, "", "not a repository")):
            result = sdk_metadata(module)
        self.assertEqual({"local": True}, result["replacement"])

    def test_git_status_failure_is_not_reported_as_clean(self):
        module = self.module()
        module["Replace"] = {"Path": "../sdk", "Dir": "/private/sdk"}
        replies = [subprocess.CompletedProcess([], 0, "abc123\n", ""),
                   subprocess.CalledProcessError(128, ["git", "status"], stderr="status failed")]
        with patch("capture.subprocess.run", side_effect=replies):
            with self.assertRaises(subprocess.CalledProcessError):
                sdk_metadata(module)


if __name__ == "__main__":
    unittest.main()
