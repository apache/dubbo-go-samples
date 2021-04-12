# Router Examples

## Background

This example will show how tag router and condition router can be used in a Dubbo go application. For more information about tag router and condition router, pls. refer to the docs:  [Conditional routing rules](https://dubbo.apache.org/en/docs/v2.7/user/examples/routing-rule/#conditional-routing-rules) and [Tag routing rules](https://dubbo.apache.org/en/docs/v2.7/user/examples/routing-rule/#tag-routing-rules).

## Config router with local config file 

Router can be configured with a local config yaml file. You can find the sample config file under "go-client/config/router_config.yml", for example:

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

You may notice '$HOST_IP' in the config file. Before run the example, you will need to replace it with the provider's IP. In order to find the provider's IP, you may consider to execute the following command:

```bash
ifconfig en0 | grep inet | grep -v inet6 | awk '{print $2}'
```

## Config router with config center

It is convenient to use environment variable 'CONF_ROUTER_FILE_PATH' to config router config in development environment. Compared to local router config file, it is more preferable to use [Dubbo Admin](https://github.com/apache/dubbo-admin) to configure router config for production use.


In order to have router configs notified, the Dubbo consumer must have a config center pre-configure like below:

```yaml
config_center:
     protocol: "zookeeper"
     address: "127.0.0.1:2181"
```