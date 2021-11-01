# Triple协议的Hessian2（Hessian+PB） 序列化方案

Triple-go支持 Pb序列化和Hessian序列化，Hessian序列化后通过嵌套入如下pb，再次进行pb序列化：
https://github.com/dubbogo/triple/blob/master/internal/codec/proto/triple_wrapper.proto

## 序列化方案选择

默认使用PB序列化，如需使用Hessian2+PB，请在配置文件中指定序列化协议为Hessian2

```yaml
# service config
services:
  "UserProvider":
    registry-ids: "demoZK"
    protocol-ids: "tri" # tri is dubbo-go3.0 protocol
    serialization: "hessian2" # hessian2 is serialization type
    interface: "org.apache.dubbo.UserProvider"
```

并按照与dubbo-go1.5.x相同的方法定义pojo、provider、consumer，即可发起调用。

## 开启服务
使用 goland 运行

triplego-hessian-client\
triplego-hessian-server