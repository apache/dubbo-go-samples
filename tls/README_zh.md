# 在Dubbo-go中使用TLS加密

## 使用方法
0.生成所需要的证书和秘钥
本示例提供已经生成好的证书和秘钥，在目录`tls/x509`下

1.配置dubbogo.yaml

客户端TLS配置：

```yaml
dubbo:
  tls_config:
    ca-cert-file: ../../../x509/server_ca_cert.pem
    tls-cert-file: ../../../x509/client2_cert.pem
    tls-key-file: ../../../x509/client2_key.pem
    tls-server-name: dubbogo.test.example.com
```

服务端TLS配置：

```yaml
dubbo:
  tls_config:
    ca-cert-file: ../../../x509/client_ca_cert.pem
    tls-cert-file: ../../../x509/server2_cert.pem
    tls-key-file: ../../../x509/server2_key.pem
    tls-server-name: dubbogo.test.example.com
```
2. 启动示例 

本示例提供了Dubbo、Grpc、Triple三种通信方式的TLS加密示例，分别位于`tls/dubbo` 、`tls/grpc` 、`tls/triple`。进入文件夹即可启动示例。

以tls/dubbo为例: 

启动服务端：

进入`tls/dubbo/go-server/cmd`,启动`server.go`

看到如下日志，则TLS配置生效

```
2022-12-01T23:39:30.690+0800    INFO    getty/getty_server.go:78        Getty Server initialized the TLSConfig configuration
```

启动客户端：

进入`tls/dubbo/go-client/cmd`，启动`client.go`

看到如下日志，则TLS配置生效

```
2022-12-01T23:40:05.998+0800    INFO    grpc/client.go:90       Grpc Client initialized the TLSConfig configuration
```