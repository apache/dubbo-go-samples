#!/bin/bash
set -euo pipefail

mvn -q -B -ntp clean compile exec:java \
  -Dexec.mainClass="org.apache.dubbo.sample.Consumer" \
  -Dexec.cleanupDaemonThreads=false
