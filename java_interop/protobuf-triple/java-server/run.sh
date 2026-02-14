#!/bin/bash
set -euo pipefail

mvn -B -ntp clean compile exec:java \
  -Dexec.mainClass="org.apache.dubbo.sample.Provider" \
  -Dexec.cleanupDaemonThreads=false
