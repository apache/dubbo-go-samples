mvn install -DSkipTests
mvn exec:java -Dexec.mainClass="org.apache.dubbo.Consumer"