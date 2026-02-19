mvn -q -e clean compile exec:java \
  -Dspotless.apply.skip=true \
  -Dspotless.check.skip=true \
  -Dexec.mainClass="org.apache.dubbo.samples.Main" \
  -Dexec.cleanupDaemonThreads=false
