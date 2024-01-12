# Triple protocol Hessian2 (Hessian+PB) serialization scheme

Triple-go supports Pb serialization and Hessian serialization. After Hessian serialization, api serialization is performed again by nesting the following api:
https://github.com/dubbogo/triple/blob/master/internal/codec/proto/triple_wrapper.proto

## Serialization scheme selection

PB serialization is used by default. To use Hessian2+PB, please specify the serialization protocol as Hessian2 in the configuration file

```yaml
# service config
services:
  "UserProvider":
    serialization: "hessian2" # hessian2 is serialization type
    interface: "org.apache.dubbo.UserProvider"
```

And define pojo, provider, and consumer in the same way as dubbo-go1.5.x, then you can initiate the call.

## Start service
Run with goland

triplego-hessian-client\
triplego-hessian-server