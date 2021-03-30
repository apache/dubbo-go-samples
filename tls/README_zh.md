## tls加密传输

dubbo-go的tls加密传输是通过dubbo-getty实现的https://github.com/apache/dubbo-getty，在使用tls之前，需要准备好

+ ca证书文件
+ client和server的密钥及证书文件

也可以使用certs文件夹中的样例文件

```shell
├── certs
│   ├── ca.key
│   ├── ca.pem
│   ├── client.key
│   ├── client.pem
│   ├── server.key
│   └── server.pem
```

### Server端

可以在运行`config.Load()`之前导入配置

```go
func init(){
  // 证书
	serverPemPath, _ := filepath.Abs("../../certs/server.pem")
  // 私钥
	serverKeyPath, _ := filepath.Abs("../../certs/server.key")
  // ca证书
	caPemPath, _ := filepath.Abs("../../certs/ca.pem")
  // 开启tls
	config.SetSslEnabled(true)
  // 导入tls配置
	config.SetServerTlsConfigBuilder(&getty.ServerTlsConfigBuilder{
		ServerKeyCertChainPath:        serverPemPath,
		ServerPrivateKeyPath:          serverKeyPath,
		ServerTrustCertCollectionPath: caPemPath,
	})
}
```

### Client端

client端的设置和server端类似，不需要设置证书

```go
func init(){
  // 私钥
	clientKeyPath, _ := filepath.Abs("../certs/ca.key")
  // ca证书
	caPemPath, _ := filepath.Abs("../certs/ca.pem")
  // 开启tls
	config.SetSslEnabled(true)
  // 导入tls配置
	config.SetClientTlsConfigBuilder(&getty.ClientTlsConfigBuilder{
		ClientPrivateKeyPath:          clientKeyPath,
		ClientTrustCertCollectionPath: caPemPath,
	})
}
```

其余设置同helloworld(https://github.com/apache/dubbo-go-samples/tree/master/helloworld)保持一致皆可。

