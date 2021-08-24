# Triple-go（Dubbo-go3.0） 示例

Triple-go 是在已有Dubbo服务治理能力的基础上，3.0 版本发布的网络协议库。

## Triple 协议介绍以及源码：

https://github.com/dubbogo/triple


## Triple-go 提供能力：
- 与 Grpc 互通，支持普通调用和流式调用
  
  [例子](./api/dubbogo-grpc/README_zh.md)
  
- 与 Triple-java 互通，目前java端支持普通调用

  [例子](./api/dubbogo-java/README_zh.md)
  

- 支持 Pb 序列化和 Hessian 序列化，Hessian 序列化后通过嵌套入如下 api，再次进行序列化：
  https://github.com/dubbogo/triple/blob/master/internal/codec/proto/triple_wrapper.proto
  
  [例子](./hessian2/README_zh.md)

