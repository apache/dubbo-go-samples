#!/bin/bash
set -euo pipefail

JDK_OPENS="--add-opens=java.base/java.lang=ALL-UNNAMED"
export MAVEN_OPTS="${MAVEN_OPTS:-} ${JDK_OPENS}"

mvn -B -ntp -e clean compile exec:java \
  -Dexec.mainClass="org.apache.dubbo.samples.Main" \
  -Dexec.cleanupDaemonThreads=false
