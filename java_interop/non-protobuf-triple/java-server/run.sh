#!/bin/bash
set -euo pipefail

mvn -q -B -ntp clean compile exec:java \
  -Dexec.mainClass="org.apache.dubbo.tri.hessian2.provider.Application" \
  -Dexec.cleanupDaemonThreads=false
