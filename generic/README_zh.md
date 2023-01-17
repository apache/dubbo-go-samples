# 泛化调用

泛化调用是在客户端没有接口信息时保证信息被正确传递的手段，即把 POJO 泛化为通用格式（如字典、字符串），一般被用于集成测试、网关等场景。

## 开始

泛化调用例子根据泛化方式划分为：

- default：默认泛化方式，即 Map 泛化方式

在每种泛化调用的例子中，又包含四种文件：

- go-server：Dubbo-Go server 例子
- go-client：Dubbo-Go client 例子
- java-client：Dubbo server 例子
- java-server：Dubbo client 例子

Dubbo Java 例子可以方便测试与 Dubbo-Go 的互通性，可以启动 java server 或 go client，或者 go server 和 java client 进行测试。

### 注册中心

本例子中使用 ZooKeeper 作为注册中心，也支持 etcd、Nacos 等作为注册中心。下面命令表示从 docker 中启动 ZooKeeper，所以需要首先确保 docker 和 docker-compose 是否已经安装。

```shell
cd ./default/go-server/docker \
  && docker-compose up -d
```
### 服务端

使用 Dubbo-Go 作为 provider，有两种方案可供选择：使用 GoLand 启动或从命令行工具启动。

使用 GoLand 启动。需要在右上角 Configurations 中选择 `v3config-generic/generic-default-go-server`，点击 Run 按钮运行即可。

从命令行工具启动。`$ProjectRootDir` 是指 dubbo-go-samples 项目根目录。

```shell
cd $ProjectRootDir/generic/default/go-server/cmd \
  && go run server.go
```

### 客户端

使用 Dubbo-Go 作为 consumer，有两种方案可供选择：使用 GoLand 启动和从命令行工具启动。

使用 GoLand 启动。需要在右上角 Configurations 中选择 `v3config-generic/generic-default-go-client`，点击 Run 按钮运行即可。

从命令行工具启动。`$ProjectRootDir` 是指 dubbo-go-samples 项目根目录。

```shell
cd $ProjectRootDir/generic/default/go-client/cmd \
  && go run client.go
```

## 将示例由接口级服务发现切换至应用级服务发现

1. 修改服务端 go-server 的配置文件，增添字段 `registry-type: service`
```
registries:
    zk:
      protocol: zookeeper
      timeout: 3s
      address: 127.0.0.1:2181
      registry-type: service
...
...
```
2. 修改客户端 go-client 的 client.go 文件

首先添加同样添加 `RegistryType: "service"`字段：
```
registryConfig := &config.RegistryConfig{
    Protocol:     "zookeeper",
    Address:      "127.0.0.1:2181",
    RegistryType: "service",
}
```
然后增添`metadataConfig`配置：
```
metadataConf := &config.MetadataReportConfig{
    Protocol: "zookeeper",
    Address:  "127.0.0.1:2181",
}
```
最后将`metadataConfig`通过`SetMetadataReport`配置进入`rootConfig`
```
...
...
rootConfig := config.NewRootConfigBuilder().
    AddRegistry("zk", registryConfig).
    SetMetadataReport(metadataConf).
    Build()
...
...
```
至此即可实现应用级服务发现的泛化调用。
