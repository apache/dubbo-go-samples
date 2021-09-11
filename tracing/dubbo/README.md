# RPC Dubbo for Dubbo-go 3.0

For api definition and go client and server startup, please refer to [dubbo-go 3.0 quickstart](https://dubbogo.github.io/zh-cn/docs/user/quickstart/3.0/quickstart.html)

## Instructions
1. Start `zookeeper` service in `docker/docker-compose.yml` file
2. Start the server side.The server startup methods of golang and java are as follows:
   1. Find the go-server folder, run the `main` function under the cmd package, and start the golang server
   2. Find the java-server folder and execute `sh run.sh` to start the java server

3. Start the client side.The client startup methods of golang and java are as follows:
   1. Find the go-client folder, run the `main` function under the cmd package, and start the golang client
   2. Find the java-client folder and execute `sh run.sh` to start the java client
4. Steps to test ziplink:
   1. Start the `zipkin` service in the `docker/docker-compose.yml` file
   2. Enable the `initZipkin()` statement of the `main` function in `go-server` and `go-client`, and commented out the `initJaeger()` statement
   3. Use the client to call the server-side service, and open http://localhost:9411/zipkin with a browser to see the tracing data
5. Steps to test jaeger:
   1. Start the `jaeger` service in the `docker/docker-compose.yml` file
   2. Enable the `initJaeger()` statement of the `main` function in `go-server` and `go-client`, and commented out the `initZipkin()` statement
   3. Use the client to call the server-side service, and open http://localhost:16686/search with a browser to see the tracing data

