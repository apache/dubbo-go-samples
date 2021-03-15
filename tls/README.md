## TLS encrypted transmission

The TLS encrypted transmission of Dubbo go is realized by Dubbo Getty https://github.com/apache/dubbo-getty . Before using TLS, you need to be prepared as follows.

+ CA certificate file
+ Key files & certificate files of client and server

You can use the sample files in certs folder.

```shell
├── certs
│   ├── ca.key
│   ├── ca.pem
│   ├── client.key
│   ├── client.pem
│   ├── server.key
│   └── server.pem
```

### Server side

Can be run in`config.Load` to import configuration before .

```go
func init(){
  //Certificate
  serverPemPath, _ := filepath.Abs ("../certs/ server.pem ")
  //Private key
  serverKeyPath, _ := filepath.Abs ("../certs/ server.key ")
  //CA certificate
  caPemPath, _ := filepath.Abs ("../certs/ ca.pem ")
  //Turn on TLS
  config.SetSslEnabled (true)
  //Import TLS configuration
  config.SetServerTlsConfigBuilder (&amp; getty.ServerTlsConfigBuilder {
    ServerKeyCertChainPath:        serverPemPath,
    ServerPrivateKeyPath:          serverKeyPath,
    ServerTrustCertCollectionPath: caPemPath,
  })
}
```
### Client side

The setting of the client side is similar to that of the server side, and there is no need to set the certificate

```go
func init(){
//Private key
	clientKeyPath, _ := filepath.Abs ("../../certs/ ca.key ")
//CA certificate
	caPemPath, _ := filepath.Abs ("../../certs/ ca.pem ")
//Turn on TLS
	config.SetSslEnabled (true)
//Import TLS configuration
  config.SetClientTlsConfigBuilder (&amp; getty.ClientTlsConfigBuilder {
    ClientPrivateKeyPath:          clientKeyPath,
		ClientTrustCertCollectionPath: caPemPath,
  })
}
```
Other settings are consistent with HelloWorld.