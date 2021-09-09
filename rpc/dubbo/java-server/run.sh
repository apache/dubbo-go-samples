mvn clean install -DSkipTests
mvn -e exec:java -Dexec.mainClass="java.org.apache.dubbo.Provider"