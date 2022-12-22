# 在Dubbo-go中使用TLS加密

## 背景

Dubbo-go在Getty/Triple/Grpc三个通信层面支持TLS链路安全通信。

## 原理

### 证书机制： 

ps: 可以先提前了解非对称加密机制。

CA(Certification Authority)负责生成根证书、签发证书等等。CA自签证书的过程如下： 

1. CA生成公钥 ca_KeyPub 和私钥 ca_KeyPri，以及基本信息表 ca_Info。ca_Info 中一般包含了CA的名称、证书的有效期等信息。
2. CA对（ca_KeyPub + ca_Info）进行散列运算，得到散列值 ca_Hash。
3. CA使用其私钥 ca_KeyPri 对 ca_Hash 进行非对称加密，得到加密的散列值 enc_ca_Hash。
4. CA将（ca_KeyPub + ca_Info + enc_ca_Hash）组合生成自签名的数字证书「ca_Cert」。这张证书称之为根证书。
5. 根证书（ca_Cert）包含的内容：ca_KeyPub + ca_Info + enc_ca_Hash。
（ca_Cert）可用于签署下一级的证书。根证书是自签名的，不需要其他机构认证。公钥私钥的生成可以利用OpenSSL等工具。

当需要签发证书时，在本地生成公钥和私钥，向CA发起CSR请求(Certificate Signing Request), CA校验此请求之后，颁发证书。过程如下：  

1. 证书申请者(S)在本地生成公钥 s_KeyPub 和私钥 s_KeyPri，以及基本信息表 s_Info。s_Info 中一般包含了证书申请者的名称、证书有效期等信息。
2. 证书申请者将 s_KeyPub、s_Info 发送给认证机构CA,即发起CSR请求。
3. CA通过某种方式验证申请者的身份之后，再加上根认证机构自己的一些信息 ca_Info，然后对它们（s_KeyPub + s_Info + ca_Info）进行散列运算，得到散列值 s_Hash。
4. CA使用其私钥 ca_KeyPri 对 s_Hash 进行非对称加密，得到加密的散列值 enc_s_Hash。
5. CA将（s_KeyPub + s_Info + ca_Info + enc_s_Hash）组合签署成数字证书（s_Cert）并发送给申请者。

申请者的证书（s_Cert）包含的内容为：s_KeyPub + s_Info + ca_Info + enc_s_Hash。

证书校验： 

用CA的公钥对 enc_s_Hash进行解密，得到s_Hash，对比其是否与 hash(s_KeyPub + s_Info + ca_Info)一致。 

### TLS机制： 
TLS 的前身是 SSL，用于通信加密，其主要过程如下：  
在最简单的TLS机制中，只需要对客户端对服务端进行校验，所以只需要服务端有证书，客户端不需要，比如我们熟悉的HTTPS。其认证过程如下：  
1. 客户端连接到服务端
2. 服务器出示其 TLS 证书(s_cert),包括服务端公钥、CA的信息等。
3. 客户端验证服务器的证书，用CA的公钥即可校验。
4. 客户端和服务器通过加密的 TLS 连接交换信息


除此之外，还有更安全的加密方式mTLS。在 mTLS 中，客户端和服务器都有一个证书，并且双方都使用它们的公钥/私钥对进行身份验证。其认证过程如下： 
1. 客户端连接到服务端
2. 服务端出示其 TLS 证书(s_cert)，包括服务端公钥、服务端CA的信息等。
3. 客户端验证服务端的证书，用服务端CA的公钥校验。
4. 客户端出示其 TLS 证书(c_cert)，包括客户端公钥、客户端CA的信息等。这里注意客户端和服务端的CA可能不是同一个。
5. 服务端验证客户端的证书，用客户端CA的公钥校验。
6. 服务端授予访问权限
7. 客户端和服务器通过加密的 TLS 连接交换信息

## 在Dubbo-go中使用TLS加密

0.生成所需要的证书和秘钥 本示例提供已经生成好的证书和秘钥，在目录`tls/x509`下

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

## 参考
https://zhuanlan.zhihu.com/p/36832100
https://www.cloudflare.com/zh-cn/learning/access-management/what-is-mutual-tls/