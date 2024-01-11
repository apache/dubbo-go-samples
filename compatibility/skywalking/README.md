# Skywalking with Dubbo-go 3.0

For api definition and go client and server startup, please refer to [dubbo-go 3.0 quickstart](https://dubbo.apache.org/zh/docs/quick-start/)

## Instructions
1. Start nacos

2. Start the server

Use goland to start skywalking/go-server/cmd/server. Note that you need to change `YOUR_SKYWALKING_DOMAIN_NAME_OR_IP` to your real environment's configuration in `server.go`, and modify the skywalking/go-server/conf/`dubbogo.yml`. 

3. Start the client

Use goland to start skywalking/go-client/cmd/client. Note that you need to change `YOUR_SKYWALKING_DOMAIN_NAME_OR_IP` to your real environment's configuration in `client.go`, and modify the skywalking/go-client/conf/`dubbogo.yml`.

