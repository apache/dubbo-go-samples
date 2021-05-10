# 路由示例

## 背景

本示例展示了如何在一个 Dubbo go 应用中使用标签路由和条件路由。有关标签路由和条件路由更多的信息，请参考文档：[条件路由](https://dubbo.apache.org/zh/docs/v2.7/user/examples/routing-rule/#%E6%9D%A1%E4%BB%B6%E8%B7%AF%E7%94%B1) and [标签路由](https://dubbo.apache.org/zh/docs/v2.7/user/examples/routing-rule/#%E6%A0%87%E7%AD%BE%E8%B7%AF%E7%94%B1%E8%A7%84%E5%88%99)。

## 使用本地配置文件配置路由

可以通过本地的一个 yaml 配置文件来配置路由。在子目录 "go-client/config/router_config.yml" 下可以找到相关的示例配置，例如：

```yaml
routerRules:
  - scope: application
    force: true
    runtime: false
    enabled: true
    valid: true
    priority: 1
    key: demo-provider
    tags:
      - name: beijing
        addresses: [192.168.1.1, $HOST_IP]
      - name: shenzhen
        addresses: [192.168.1.3, 192.168.1.4]
```

在这个配置文件中您可能注意到了变量 '$HOST_IP'。运行示例之前，您需要将它替换成您的服务提供方的 IP，可以通过运行以下命令的方式获取：

```bash
ifconfig en0 | grep inet | grep -v inet6 | awk '{print $2}'
```

## 通过配置中心配置路由

通过环境变量 'CONF_ROUTER_FILE_PATH' 指定本地路由配置的方式对于开发环境来说是十分方便的。对比于这种方式，在生产上使用时，更建议通过 [Dubbo Admin](https://github.com/apache/dubbo-admin) 来下发路由配置。

为了让路由配置能够被 Dubbo go 应用感知到，应用需要提前配置好配置中心，如下所示：

```yaml
config_center:
     protocol: "zookeeper"
     address: "127.0.0.1:2181"
```