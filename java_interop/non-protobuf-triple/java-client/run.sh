#!/bin/bash
set -euo pipefail

mvn -q -B -ntp clean compile exec:java \
  -Dexec.mainClass="org.apache.dubbo.tri.hessian2.client.Application" \
  -Dexec.cleanupDaemonThreads=false
