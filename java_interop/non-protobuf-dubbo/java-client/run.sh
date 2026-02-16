#!/bin/bash
set -euo pipefail

mvn -B -ntp clean compile exec:java \
  -Dexec.mainClass="org.apache.dubbo.hessian2.client.Application" \
  -Dexec.cleanupDaemonThreads=false
