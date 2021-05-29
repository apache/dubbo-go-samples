# Tengine Calls Dubbo-go 示例

## 背景

Tengine 是由淘宝网在Nginx的基础上发起的Web服务器项目。针对大访问量网站的需求，添加了很多高级功能和特性，比如 `2.3.2` 开始支持 Dubbo 协议，本示例将演示基于 Dubbo 协议，Tengine 调用 Dubbo-go 服务。

## 介绍

```markdown
.
├── README.md
├── README_zh.md
└── go-server
    ├── cmd
    ├── conf
    ├── docker
    └── pkg
```

- go-server ：服务提供者
- docker ： 提供 Tengine 与 Zookeeper 

## 安装 Tengine

### 使用 Docker 

使用 Docker :

直接启动 [/dubbo-go-sample/tengine/docker/docker-compose.yml](docker-compose.yml)

本示例使用最新版本 tengine `v2.3.3`，事实上从 2.3.2 开始支持 `dubbo` 协议

### 手动安装

参考 [tengine](https://github.com/alibaba/tengine) 

#### Clone Tengine
```
git clone https://github.com/alibaba/tengine.git
```
#### 安装其他支持库
```
cd ./tengine

wget https://ftp.pcre.org/pub/pcre/pcre-8.43.tar.gz
tar xvf pcre-8.43.tar.gz

wget https://www.openssl.org/source/openssl-1.0.2s.tar.gz
tar xvf openssl-1.0.2s.tar.gz

wget http://www.zlib.net/zlib-1.2.11.tar.gz
tar xvf zlib-1.2.11.tar.gz
```

### 构建 Tengine
```
./configure --add-module=./modules/mod_dubbo --add-module=./modules/ngx_multi_upstream_module --add-module=./modules/mod_config --with-pcre=./pcre-8.43/ --with-openssl=./openssl-1.0.2s/ --with-zlib=./zlib-1.2.11
make
sudo make install
```

### 修改配置

修改 tengine 配置文件 `/usr/local/nginx/conf/nginx.conf` 

```
worker_processes  1;

events {
    worker_connections  1024;
}


http {
    include       mime.types;
    default_type  application/octet-stream;

    sendfile        on;

    server {
        listen       8080;
        server_name  localhost;
        
        #pass the Dubbo to Dubbo Provider server listening on 127.0.0.1:20880
        location / {
            dubbo_pass_all_headers on;
            dubbo_pass_set args $args;
            dubbo_pass_set uri $uri;
            dubbo_pass_set method $request_method;
        
            dubbo_pass org.apache.dubbo.UserProvider 0.0.0 GetUser dubbo_backend;
        }
    }

    #pass the Dubbo to Dubbo Provider server listening on 127.0.0.1:20880
    upstream dubbo_backend {
        multi 1;
        server 127.0.0.1:20000;
    }
}
```

### 启动 Tengine

```
/usr/local/nginx/sbin/nginx
```

重启、停止
```
#restart
/usr/local/nginx/sbin/nginx -s reload
#stop
/usr/local/nginx/sbin/nginx -s stop
```

## 启动服务提供者

注意：
配置环境变量，指定服务加载所需配置文件路径

go-server:

```shell
APP_LOG_CONF_FILE=direct/go-server/conf/log.yml;
CONF_PROVIDER_FILE_PATH=direct/go-server/conf/server.yml
```

具体请参阅根目录中的 [HOWTO.md](../HOWTO_zh.md) 来运行本例。

## 查看效果

执行：

```shell
curl http://127.0.0.1:8080/dubbo -i
```

响应：

```shell
HTTP/1.1 200 OK
Server: Tengine/2.3.3
Date: Sat, 29 May 2021 15:41:30 GMT
Content-Type: application/octet-stream
Transfer-Encoding: chunked
Connection: keep-alive
Age: 18
ID: A001
Name: Alex Stocks
Time: 1622302890448
```

一个简单的 Tengine Calls Dubbo-go 的例子已经完成了

更多功能请参考：

[tengine/../ngx_http_dubbo_module_cn](https://github.com/alibaba/tengine/blob/master/docs/modules/ngx_http_dubbo_module_cn.md)




