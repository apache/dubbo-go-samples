# Mesh Route for Dubbo-go 3.0

For api definition and go client and server startup, please refer to [dubbo-go 3.0 quickstart](https://dubbogo.github.io/zh-cn/docs/user/quickstart/3.0/quickstart.html)

## Instructions
1. Start `zookeeper` service by docker with `integrate_test/dockercompose/docker-compose.yml` or executable binary file.
2. Start the server service.
3. Get dynamic configuration and publish mesh route config.
4. Use the client to call the server-side service.