mvn install -DSkipTests
mvn exec:java -Dexec.mainClass="com.apache.dubbo.sample.basic.ApiConsumer"