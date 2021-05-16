# Shopping-Order Example

## 背景

电商交易等场景，业务对分布式事务有很高的要求......

## 介绍

一个简单的电商交易场景，基于 [Dubbo-Go](https://github.com/apache/dubbo-go) 和 [Seata-Golang](https://github.com/opentrx/seata-golang)  **AT事务模式** 实现的分布式事务解决方案示例， 
关于 Seata 提供的其他事务模式，详细请参阅 [Seata](https://github.com/seata/seata) 


### 模块

```shell
.
├── go-client
│   ├── cmd
│   ├── conf
│   ├── pkg
│   └── profiles
│       └── dev
├── go-server-common
│   └── filter
├── go-server-order
│   ├── cmd
│   ├── conf
│   ├── docker
│   └── pkg
│       └── dao
└── go-server-product
    ├── cmd
    ├── conf
    ├── docker
    └── pkg
        └── dao

```
- go-client : 服务消费者
- go-server-common : 项目公共模块，包含过滤器等
- go-server-order : 订单服务提供者
- go-server-product : 库存服务提供者

### 运行

请参阅根目录中的 [HOWTO.md](../HOWTO_zh.md) 来运行本例。

- [ ] Go
- [ ] Zookeeper or Nacos ...
- [ ] Mysql
- [ ] **Seata-Golang Server (Transaction Coordinator)**


#### 1. 环境配置

配置环境变量，指定服务加载所需配置文件路径

1.1 go-server-order

1.1.1 环境配置

```shell
APP_LOG_CONF_FILE=shopping-order/go-server-order/conf/log.yml;
CONF_PROVIDER_FILE_PATH=shopping-order/go-server-order/conf/server.yml;
SEATA_CONF_FILE=shopping-order/go-server-order/conf/seata.yml
```

1.1.2 配置文件

修改 [server.yml](go-server-order/conf/server.yml)

- 注册中心及地址

修改 [seata.yml](go-server-order/conf/seata.yml)

- Seata AT DNS

```shell
at:
  dsn: "username:password@tcp(mysql:3306)/seata_order?timeout=1s&readTimeout=1s&writeTimeout=1s&parseTime=true&loc=Local&charset=utf8mb4,utf8"
  report_retry_count: 5
  report_success_enable: false
  lock_retry_interval: 10ms
  lock_retry_times: 30
```

1.2 go-server-product:

1.2.1 环境配置：

```shell
APP_LOG_CONF_FILE=shopping-order/go-server-product/conf/log.yml;
CONF_PROVIDER_FILE_PATH=shopping-order/go-server-product/conf/server.yml;
SEATA_CONF_FILE=shopping-order/go-server-product/conf/seata.yml
```

1.2.2 配置文件

修改 [server.yml](go-server-product/conf/server.yml)

- 注册中心及地址

修改 [seata.yml](go-server-product/conf/seata.yml)

- Seata AT DNS

```shell
at:
  dsn: "username:password@tcp(mysql:3306)/seata_product?timeout=1s&readTimeout=1s&writeTimeout=1s&parseTime=true&loc=Local&charset=utf8mb4,utf8"
  report_retry_count: 5
  report_success_enable: false
  lock_retry_interval: 10ms
  lock_retry_times: 30
```

1.3 go-client:

环境配置：

```shell
APP_LOG_CONF_FILE=shopping-order/go-client/conf/log.yml;
CONF_CONSUMER_FILE_PATH=shopping-order/go-client/conf/client.yml;
SEATA_CONF_FILE=shopping-order/go-client/conf/seata.yml
```

#### 2. 启动注册中心
   
本示例使用 Zookeeper 做注册中心， 可以直接运行 docker zookeeper 环境，配置详情请参阅 `docker-compose.yml`

#### 3. 启动TC(Transaction Coordinator)事务协调者

详情可参考：
[Seata-Golang](https://github.com/opentrx/seata-golang) 社区，
[Seata](https://github.com/seata/seata) 社区

#### 4. 启动服务提供者

> - 订单服务
> - 库存服务

`${project_dir}/.run/shopping-order`

```shell
  .
├── go-client.run.xml
├── order-go-server.run.xml
└── product-go-server.run.xml
```

#### 5. 启动消费者

---
请参阅根目录中的 [HOWTO.md](../HOWTO_zh.md) 来运行本例。
