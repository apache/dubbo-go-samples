# Dubbo Go registry-type support all

| 参数         | 注册方式       |
|------------|------------|
| interface  | 接口级        |
| service    | 应用级        |
| all        | 应用级 && 接口级 |


## 使用方法
1. 启动 zk/nacos

2. 启动服务端
   将启动配置修改为： registry-type: all
```shell
   cd registry/all/go-server
   DUBBO_GO_CONFIG_PATH=registry/all/nacos/go-server/conf/dubbogo.yml && go run server.go
   ```

3. 启动客户端
   由于服务注册方式支持双注册，客户端使用 service/interface 两种参数均可访问服务端
```shell
  cd registry/all/go-client
  DUBBO_GO_CONFIG_PATH=registry/all/nacos/go-client/conf/dubbogo.yml && go run client.go 
```




