# Triple-go (Dubbo-go3.0) example

Triple-go is a network protocol library released in version 3.0 based on the existing Dubbo service management capabilities.

## Triple protocol introduction and source code:

https://github.com/dubbogo/triple


## Triple-go provides capabilities:
-Intercommunication with Grpc, supporting normal calls and streaming calls

[Example](./pb/dubbogo-grpc/README.md)

-Intercommunication with Triple-java, currently the java side supports ordinary calls

[Example](./pb/dubbogo-java/README.md)


-Support Pb serialization and Hessian serialization. After Hessian serialization, it can be serialized again by nesting the following pb:
https://github.com/dubbogo/triple/blob/master/internal/codec/proto/triple_wrapper.proto

[Example](./hessian2/README.md)

