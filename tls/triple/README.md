# Dubbo-go Triple mTLS Encrypted Communication

This document describes how to implement mTLS (mutual Transport Layer Security) encrypted communication through Dubbo-go's Triple protocol.

The certificate file is created using the script provided by grpc-go. In the x509 directory, you can run `create.sh` to create it.

## 1. Background

In distributed systems, secure encrypted communication is a key means of protecting data from unauthorized access. mTLS further enhances security through mutual authentication, as it requires both the client and server to present and validate each other's certificates.

## 2. Using Custom Codec

In this example, we use msgpack as a custom serialization/deserialization Codec. While the Triple protocol supports protojson by default, we chose msgpack because it may offer better performance in certain scenarios.

### 2.1 Why Not Use JSON Directly?

- **Performance Considerations**: In certain high-performance scenarios, msgpack has higher serialization and deserialization speeds than JSON.
- **Protocol Flexibility**: msgpack helps support flexible protocol design, occupying less bandwidth during data transmission.

### 2.2 TODO: Support JSON Codec

Although msgpack is currently in use, we plan to add support for JSON Codec to the Triple protocol in future versions. This way, users can choose the most suitable serialization format based on their needs.

## 3. Certificate File Path Issues

### 3.1 Limitations of Integration Testing

Currently, the script for launching integration tests cannot handle certificate files with relative paths, which limits the implementation and execution of automated tests.

### 3.2 Future Improvement Directions

We plan to:
- **Update Scripts**: Support relative path parsing to achieve more flexible file path configuration.
- **Increase Documentation**: Provide detailed guidelines for certificate generation and path configuration.
- **Automated Testing**: Add more integration test cases for mTLS scenarios to enhance security and stability.
