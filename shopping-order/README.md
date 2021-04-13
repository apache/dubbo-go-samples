# Shopping-Order Example

## Backend

Business scenarios such as e-commerce transactions have high requirements for distributed transactions...

## Introduction

A simple e-commerce transaction scenario, an example of a distributed transaction solution implemented based on [Dubbo-Go](https://github.com/apache/dubbo-go) and [Seata-Golang](https://github.com/opentrx/seata-golang)  **AT transaction patterns**. See Seata for details on other transaction modes offered by  [Seata](https://github.com/seata/seata).

### Modules

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
- go-client : The Service Consumer
- go-server-common : Public Module: Common components, Filter ...
- go-server-order : The Service Provider Of Order
- go-server-product : The Service Provider Of Product

### Start

Refer to  [HOWTO.md](../HOWTO_zh.md) under the root directory to run this sample.

- [ ] Go
- [ ] Zookeeper or Nacos ...
- [ ] Mysql
- [ ] **Seata-go Server (Transaction Coordinator)**


#### 1. Environment Configuration

Configure the environment variable to specify the configuration file path required for the service to load

1.1 go-server-order

1.1.1 environment configuration

```shell
APP_LOG_CONF_FILE=shopping-order/go-server-order/conf/log.yml;
CONF_PROVIDER_FILE_PATH=shopping-order/go-server-order/conf/server.yml;
SEATA_CONF_FILE=shopping-order/go-server-order/conf/seata.yml
```

1.1.2 configuration file

modify [server.yml](go-server-order/conf/server.yml)

- modify registration centre 

modify [seata.yml](go-server-order/conf/seata.yml)

- modify AT model DNS of Seata

```shell
at:
  dsn: "username:password@tcp(mysql:3306)/seata_order?timeout=1s&readTimeout=1s&writeTimeout=1s&parseTime=true&loc=Local&charset=utf8mb4,utf8"
  report_retry_count: 5
  report_success_enable: false
  lock_retry_interval: 10ms
  lock_retry_times: 30
```



1.2 go-server-product:

1.2.1 environment configuration

```shell
APP_LOG_CONF_FILE=shopping-order/go-server-product/conf/log.yml;
CONF_PROVIDER_FILE_PATH=shopping-order/go-server-product/conf/server.yml;
SEATA_CONF_FILE=shopping-order/go-server-product/conf/seata.yml
```

1.2.2 configuration file

modify [server.yml](go-server-product/conf/server.yml)

- modify registration centre

modify [seata.yml](go-server-product/conf/seata.yml)

- modify AT model DNS of Seata

```shell
at:
  dsn: "username:password@tcp(mysql:3306)/seata_product?timeout=1s&readTimeout=1s&writeTimeout=1s&parseTime=true&loc=Local&charset=utf8mb4,utf8"
  report_retry_count: 5
  report_success_enable: false
  lock_retry_interval: 10ms
  lock_retry_times: 30
```


1.3 go-client:

environment configuration:

```shell
APP_LOG_CONF_FILE=shopping-order/go-client/conf/log.yml;
CONF_CONSUMER_FILE_PATH=shopping-order/go-client/conf/client.yml;
SEATA_CONF_FILE=shopping-order/go-client/conf/seata.yml
```

#### 2. Start The Registry

This example uses ZooKeeper as the registry, so you can run the Docker ZooKeeper environment directly. See [docker-compose.yml](go-server/docker/docker-compose.yml)

#### 3. Start TC(Transaction Coordinator)

See：
[seata-golang](https://github.com/opentrx/seata-golang) ，
[seata](https://github.com/seata/seata) 

#### 4. Start The Provider

> - Order Service
> - Product services

`${project_dir}/.run/shopping-order`

```shell
  .
├── go-client.run.xml
├── order-go-server.run.xml
└── product-go-server.run.xml
```


#### 5. Start The Consumer

---

Refer to  [HOWTO.md](../HOWTO_zh.md) under the root directory to run this sample.
