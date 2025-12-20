# TLS Example

## Description

This example demonstrates how to use TLS (based on X.509 certificates) in Dubbo-Go and Dubbo-Java to enable encrypted communication and/or mutual authentication between a client and a server. More importantly, this example showcases the **cross-language interoperability between Dubbo-Go and Dubbo-Java**—through the Triple protocol and Protobuf serialization, a Go client can seamlessly call a Java server, and a Java client can call a Go server. The example includes client and server sample programs for both Go and Java, as well as scripts for generating test certificates using X.509.

## Directory Structure

* **go-client/**: Go client example program
* **go-server/**: Go server example program
* **java-server/**: Java server example program (Dubbo-Java provider)
* **java-client/**: Java client example program (Dubbo-Java consumer)
* **proto/**: Proto file and generated code for the `greet` service
* **x509/**: Scripts and example certificates for generating/storing test certificates

## Prerequisites

* Go (recommended version 1.18+) for running Go clients and servers
* Java Development Kit (JDK 8+) and Maven (3.6.0+), required only when running Java servers or clients
* On Windows, it is recommended to use Git Bash or WSL to run the certificate generation script (`x509/create.sh` uses OpenSSL).

## Generating Test Certificates

1. Navigate to the `x509` directory and run the certificate generation script:

    * On Unix-based systems, execute the following command:

      ```bash
      cd tls/x509 && ./create.sh  
      ```
    * On Windows, if OpenSSL or Bash is not installed, you can run the script using WSL/Git Bash or manually generate the certificates using OpenSSL following the configuration in `x509/openssl.cnf`.
2. The generated certificates will be stored in the `x509/` directory, including:

    * `server_ca_*.pem`
    * `client_ca_*.pem`
    * `server{1,2}_*.pem`
    * `client{1,2}_*.pem`

## Running the Example

### Option 1: Go Server + Go Client

#### 1. Start the Go Server

In the tls directory, run the following command to start the Go server:

```bash
go run ./go-server/cmd  
```

The server will load the server certificates and CA from the `x509/` directory and listen on the address specified in the configuration. If you need to customize the configuration, please modify the `server` program or the relevant parts in the source code.

#### 2. Start the Go Client

In another terminal in the tls directory, run the following command to start the Go client:

```bash
go run ./go-client/cmd  
```

The client will use the certificates from the `x509/` directory to establish a TLS connection with the server and invoke the `greet` service.

### Option 2: Java Server + Go Client (Dubbo-Go ↔ Dubbo-Java Interoperability)

This demonstrates **interoperability between Dubbo-Go and Dubbo-Java** using the Triple protocol over TLS.

#### 1. Start the Java Server (Dubbo-Java Provider)

In the tls directory, navigate to the `java-server` subdirectory and start the Maven project:

```bash
cd ./java-server
mvn clean compile
mvn exec:java -Dexec.mainClass="org.apache.dubbo.samples.tls.provider.TlsTriProvider"
```

The Java server will start on port 20000 using TLS with the certificates from `x509/server2_cert.pem` and `x509/server2_key.pem`.

#### 2. Start the Go Client

In another terminal, in the tls directory, run the Go client:

```bash
go run ./client/cmd
```

The Go client will connect to the Java server via TLS and invoke the `greet` service. You should see output like:

```
Greet response: hello world from Java provider
```

### Option 3: Go Server + Java Client

This option demonstrates the reverse interoperability—a Java client calling a Go server.

#### 1. Start the Go Server

In the tls directory, run the following command to start the Go server:

```bash
go run ./go-server/cmd
```

The Go server will start on port 20000 with TLS encryption enabled.

#### 2. Start the Java Client

In another terminal, in the tls directory, navigate to the `java-client` subdirectory and start the Maven project:

```bash
cd ./java-client
mvn clean compile
mvn exec:java -Dexec.mainClass="org.apache.dubbo.samples.tls.consumer.TlsTriProtoConsumer"
# To customize the target host and TLS authority (SNI), add: -Dtls.host=127.0.0.1 -Dtls.authority=dubbogo.test.example.com
```

The Java client will connect to the Go server via TLS and invoke the `greet` service. You should see output like:

```
Greet response: hello world
```

## Dubbo-Go and Dubbo-Java Interoperability

This example demonstrates **cross-language interoperability** between Dubbo-Go and Dubbo-Java frameworks:

* **Protocol**: Both use the Triple protocol (based on gRPC/HTTP2)
* **Serialization**: Protobuf ensures language-agnostic data exchange
* **TLS/SSL**: Both support X.509 certificate-based encryption and authentication
* **Service Interface**: Defined in `proto/greet.proto`, code is generated for both Go and Java

## Notes

* The certificate paths and whether mutual authentication is enabled depend on the files loaded by the example programs. Please check `tls/go-server/cmd/main.go` and `tls/go-client/cmd/main.go` to understand the specific behavior and available command-line parameters.
* On Windows, running the `create.sh` script may require WSL/Git Bash or manually running OpenSSL commands.
* This example is intended for educational and testing purposes only. The example certificates should not be used in a production environment.
