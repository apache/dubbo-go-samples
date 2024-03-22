# dubbogo-java

## Contents
- protobuf: Using the struct definitions from proto files
- server
- client
Please note that this sample is coded using dubbo-go 3.2.0-rc1.
The combinations we have tested include the following:

- [x] java-client communicating with a dubbogo-server
- [x] java-server communicating with a dubbogo-client

## Generating Code
- Java
  - 1.Add build to pom.xml
    ```xml
      <build>
        <extensions>
            <extension>
                <groupId>kr.motd.maven</groupId>
                <artifactId>os-maven-plugin</artifactId>
                <version>1.6.1</version>
            </extension>
        </extensions>
        <plugins>
            <plugin>
                <groupId>org.xolstice.maven.plugins</groupId>
                <artifactId>protobuf-maven-plugin</artifactId>
                <version>0.6.1</version>
                <configuration>
                    <protocArtifact>com.google.protobuf:protoc:3.19.4:exe:${os.detected.classifier}</protocArtifact>
                    <outputDirectory>${project.basedir}/../build/protobuf/java</outputDirectory>
                    <protocPlugins>
                        <protocPlugin>
                            <id>dubbo</id>
                            <groupId>org.apache.dubbo</groupId>
                            <artifactId>dubbo-compiler</artifactId>
                            <version>${dubbo.version}</version>
                            <mainClass>org.apache.dubbo.gen.tri.Dubbo3TripleGenerator</mainClass>
                        </protocPlugin>
                    </protocPlugins>
                </configuration>
                <executions>
                    <execution>
                        <goals>
                            <goal>compile</goal>
                        </goals>
                    </execution>
                </executions>
            </plugin>
        </plugins>
    </build>
    ```
  - 2.Generate using mvn
  ```shell
  mvn clean install
  ```
- go
  - Generate code with protoc
    ```shell
    protoc --go_out=. --go_opt=paths=source_relative --go-triple_out=. greet.proto 
    ```


## Running the Application
1. Start the server:
    - Use goland to start triple/gojava-go-server
    - Execute `sh run.sh` in the java-server folder to start the java server
2. Start the client
    - Use goland to start triple/gojava-go-client
    - Execute `sh run.sh` under the java-client folder to start the java client

## Notes
1. Interface naming must be consistent
   - java-server: GreeterImpl
   - go-client: The configuration should be similarly defined as follows
   ```yml
     Consumer:
       services:
         GreeterConsumer:
           # interface is for registry
           interface: org.apache.dubbo.sample.GreeterImpl
   ```