# Use TLS encryption in Dubbo go

## Usage

0. Generate the required certificate and secret key

This example provides the generated certificate and secret key under the directory `tls/x509`

1. Configure dubbogo.yaml

Client TLS configuration:

```yaml
dubbo:
   tls_config:
      ca-cert-file: ../../../x509/server_ca_cert.pem
      tls-cert-file: ../../../x509/client2_cert.pem
      tls-key-file: ../../../x509/client2_key.pem
      tls-server-name: dubbogo.test.example.com
```

Server TLS configuration:

```yaml
dubbo:
   tls_config:
      ca-cert-file: ../../../x509/client_ca_cert.pem
      tls-cert-file: ../../../x509/server2_cert.pem
      tls-key-file: ../../../x509/server2_key.pem
      tls-server-name: dubbogo.test.example.com
```

2. Startup example

This example provides TLS encryption examples of Dubbo, Grpc and Triple communication modes, respectively located in

`tls/dubbo` 、`tls/grpc` 、`tls/triple`。 Enter the folder to launch the sample.

Take tls/dubbo as an example:

* step1: Start the server:

Enter 'tls/dubbo/go server/cmd' and start 'server.go`

The TLS configuration takes effect when you see the following logs

```
2022-12-01T23:39:30.690+0800 INFO getty/getty_ server. go:78 Getty Server initialized the TLSConfig configuration
```

* step2: Start client:

Enter 'tls/dubbo/go client/cmd' and start 'client.go`

The TLS configuration takes effect when you see the following logs

```
2022-12-01T23:40:05.998+0800 INFO grpc/client. go:90 Grpc Client initialized the TLSConfig configuration
```