# TLS Example

## Description

This example demonstrates how to use TLS (based on X.509 certificates) in Dubbo-Go to enable encrypted communication and/or mutual authentication between a client and a server. The example includes a simple `greet` service, client and server sample programs, and scripts for generating test certificates using X.509.

## Directory Structure

* **client/**: Client example program
* **server/**: Server example program
* **proto/**: Proto file and generated code for the `greet` service
* **x509/**: Scripts and example certificates for generating/storing test certificates

## Prerequisites

* Go (recommended version 1.18+)
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

### 1. Start the Server

In the root project directory, run the following command to start the server:

```bash
go run ./tls/server/cmd  
```

The server will load the server certificates and CA from the `x509/` directory and listen on the address specified in the configuration. To customize this, modify the `server` program or the source code as needed.

### 2. Start the Client

In another terminal, run the following command to start the client:

```bash
go run ./tls/client/cmd  
```

The client will use the client certificates and CA from the `x509/` directory to establish a TLS connection with the server and invoke the `greet` service.

## Notes

* The certificate paths and whether mutual authentication is enabled depend on the files loaded by the example programs. Please check `tls/server/cmd/main.go` and `tls/client/cmd/main.go` to understand the specific behavior and available command-line parameters.
* On Windows, running the `create.sh` script may require WSL/Git Bash or manually running OpenSSL commands.
* This example is intended for educational and testing purposes only. The example certificates should not be used in a production environment.
