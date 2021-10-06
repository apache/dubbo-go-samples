mvn clean install -DSkipTests
mvn -e exec:java -Dexec.mainClass="com.apache.dubbo.sample.basic.ApiConsumer"