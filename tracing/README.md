# RPC Dubbo for Dubbo-go 3.0

If u wanna know how to startup the client/server and the definition of APIs, please refer to [dubbo-go 3.0 quickstart](https://dubbogo.github.io/zh-cn/docs/user/quickstart/3.0/quickstart.html)

## Instructions
1. Start `zookeeper` service in `docker/docker-compose.yml` file
2. Start the server side. 
   1. Find the go-server folder, run the `main` function under the cmd package, and start the golang server
3. Start the client side. 
   1. Find the go-client folder, run the `main` function under the cmd package, and start the golang client
4. Steps to test ziplink:
   1. Start the `zipkin` service in the `docker/docker-compose.yml` file
   2. Enable the `initZipkin()` statement of the `main` function in `go-server` and `go-client`, and commented out the `initJaeger()` statement
   3. Use the client to call the server-side service, and open http://localhost:9411/zipkin with a browser to see the tracing data
5. Steps to test jaeger:
   1. Start the `jaeger` service in the `docker/docker-compose.yml` file
   2. Enable the `initJaeger()` statement of the `main` function in `go-server` and `go-client`, and commented out the `initZipkin()` statement
   3. Use the client to call the server-side service, and open http://localhost:16686/search with a browser to see the tracing data

