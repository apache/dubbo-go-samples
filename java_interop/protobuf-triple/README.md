# dubbogo-java

## Contents
- protobuf: Using the struct definitions from proto files
- server
- client
Please note that this sample is coded using dubbo-go 3.2.0-rc1.
The combinations we have tested include the following:

- [x] java-client communicating with a dubbogo-server
- [x] java-server communicating with a dubbogo-client

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