## 配置API

配置API旨在支持用户通过API的形式来使用dubbo-go框架，而无需定义配置文件。
配置API具有更高的灵活性，可以在代码中动态地写入配置条目，从而启动框架服务，或者生成自己需要使用的实例。例如注册中心、配置中心等模块。

### 目录结构

-configcenter/

使用配置API，从远程配置中心读取框架启动配置，详情见 [README](rpc/README.md)

- rpc/

使用配置API启动provider端和consumer端，发起调用

- subModule/

使用配置API 获取框架特定组件的实例
