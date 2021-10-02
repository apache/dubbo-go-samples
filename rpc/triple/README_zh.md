# Triple 示例

Triple-go 3.0 版本的网络协议库。获取更多信息请查看 [dubbogo/triple](https://github.com/dubbogo/triple) 仓库。Triple 协议是在已有 Dubbogo 框架的基础上的扩展的3.0新网络协议，支持了pb序列化，可与 Dubbo3.0 互通、与 gRPC 互通，支持普通 RPC 调用与流式 RPC 调用等功能，是 Dubbo 生态在云原生时代的主推协议。

## Samples内容

- [codec-extension](./codec-extension): 用户自定义序列化方式例子
- [hessian2](./hessian2): Hessian2 序列化方式例子
- [msgpack](./msgpack): Msgpack 序列化方式例子
- [pb](./pb): 使用 ProtoBuf(PB) 序列化方案的例子
  - [dubbogo-grpc](./pb/dubbogo-grpc): Triple 和 gRPC 互通案例
  - [dubbogo-java](./pb/dubbogo-java): Triple-java 和 Triple-go 互通案例

## 如何配置

- 服务端

```yaml
dubbo:
	protocols: # 框架协议配置
		myProtocol: # 自定义一个协议 Key
			name: tri # 协议名，支持tri/dubbo/grpc/jsonrpc
			port: 20000 # 暴露端口
			
	provider: 
		services: 
			MyProvider: # 服务提供者结构类名
			 protocol: myProtocol # 自定义的协议 Key，与上方 myProtocol 对应
			 interface: org.apache.dubbogo.MyProvider # 用户自定义的接口名
		#  serialization: hessian2 可选字段，可以指定序列化类型：pb/hessian2/自定义
		#	 默认使用 pb 序列化

```

- 客户端

```yml
dubbo:
  consumer:
    references:
      ClientImpl: # 客户端结构类名
        protocol: tri # 协议名，支持tri/dubbo/grpc/jsonrpc，需与服务端对应
        interface: org.apache.dubbo.demo.Greeter # 用户自定义的接口名
     #  serialization: hessian2 可选字段，可以指定序列化类型：pb/hessian2/自定义
		 #	默认使用 pb 序列化，需要与服务端对应
```

## 运行示例：

以 pb/dubbogo-grpc 下的 dubbogo-client 调用 dubbogo-server为例

启动zk，监听127.0.0.1:2181端口。如本机已安装docker，可以直接执行下面的命令来启动所有运行samples的依赖组件：zk(2181), nacos(8848), etcd(2379)。

`docker-compose -f {PATH_TO_SAMPLES_PROJECT}/integrate_test/dockercompose/docker-compose.yml up -d`

### 通过 Goland 运行

![](../../.images/samples-rpc-triple-server.png)

服务端启动完毕后，启动客户端

![](../../.images/samples-rpc-triple-client.png)



### 通过命令行运行

- 服务端

`cd rpc/triple/pb/dubbogo-grpc/go-server/cmd` # 进入仓库目录

`export DUBBO_GO_CONFIG_PATH="../conf/dubbogo.yml`# 设置配置文件环境变量

`go run .` # 启动服务

- 客户端

`cd rpc/triple/pb/dubbogo-grpc/go-server/cmd` # 进入仓库目录

`export DUBBO_GO_CONFIG_PATH="../conf/dubbogo.yml`# 设置配置文件环境变量

`go run .` # 启动客户端发起调用



### 调用成功

通过上述任一方式启动，可看到客户端打印如下信息，调用成功：

`INFO    cmd/client.go:108       Receive user = name:"Hello laurence" id:"12345" age:21`
