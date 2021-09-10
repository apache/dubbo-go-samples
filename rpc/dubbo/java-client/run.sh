mvn clean install -DSkipTests
mvn -e exec:java -Dexec.mainClass="org.apache.dubbo.Consumer"