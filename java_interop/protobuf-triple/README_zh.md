# dubbogo-java

使用同一个proto文件实现dubbo的java和go互通
## Contents

- protobuf: 使用 proto 文件的结构体定义
- server
- client

请注意，该样例使用dubbo-go 3.2.0-rc1编写
我们测试的组合包括:

- [x] java-client -> dubbogo-server
- [x] java-server -> dubbogo-client
## 生成 code
- java
  - 1.在pom.xml中添加build
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
  - 2.使用mvn生成
  ```shell
  mvn clean install
  ```
- go
  - 使用protoc生成code
  ```shell
  protoc --go_out=. --go_opt=paths=source_relative --go-triple_out=. greet.proto 
  ```
  
## 运行
1. 启动服务端
   - 使用 goland 启动 triple/gojava-go-server
   - 在 java-server 文件夹下执行 `sh run.sh` 启动 java server
2. 启动客户端
   - 使用 goland 启动 triple/gojava-go-client
   - 在 java-client 文件夹下执行 `sh run.sh` 启动 java client

## 注意
1. 接口命名须一致
   - java-server: GreeterImpl
   - go-client: 在conf中应类似如下定义
   ```yml
     Consumer:
       services:
         GreeterConsumer:
           # interface is for registry
           interface: org.apache.dubbo.sample.GreeterImpl
   ```
      