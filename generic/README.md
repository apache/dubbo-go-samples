# Generic Invocation

Generic invocation is mainly used when the client does not have API interface or model class, all POJOs in parameters and return values are represented by map or other generic data structures. Commonly used for framework integration such as: implementing a common service testing framework, all service implementations can be invoked via `GenericService`. For more information please visit our documentation.

## Dubbo Protocol

### Instructions

1. Start zookeeper

   ```shell
   cd ./default/go-server/docker \
     && docker-compose up -d
   ```

2. Start the server to run provider.

   1. go

      Use goland to start generic-dubbo-go-server

   2. java

      Use goland to start generic-dubbo-java-server
      
      or

      Execute `sh run.sh` in the java-server folder to start the java server

3. Start the client to run consumers to initiate generic invocation.

   1. go

      Use goland to start generic-dubbo-go-client

   2. java

      Use goland to start generic-dubbo-java-client
   
      or

      Execute `sh run.sh` in the java-client folder to start the java client