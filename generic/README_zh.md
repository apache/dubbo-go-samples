# 泛化调用

泛化接口调用方式主要用于客户端没有 API 接口及模型类元的情况，参数及返回值中的所有 POJO 均用 Map 表示，通常用于框架集成，比如：实现一个通用的服务测试框架，可通过`GenericService`调用所有服务实现。更多信息请参与文档。

## 开始

你需要提供提供者（provider）和消费者（consumer）的配置文件信息，在这个例子中，使用`default`文件夹作为示例。

```shell
export CONF_PROVIDER_FILE_PATH="$(pwd)/default/go-server/conf/server.yml"
export CONF_CONSUMER_FILE_PATH="$(pwd)/default/go-client/conf/client.yml"
```

ZooKeeper也是必须的，你可以使用docker-compose启动它。

```shell
cd ./default/go-server/docker \
  && docker-compose up -d
```

### Map：默认方式

Map的例子放在`default`文件夹中。通过下面的代码启动提供者。

```shell
cd ./default/go-server/cmd \
  && go run server.go
```

运行消费者来发起泛化调用。

```shell
cd ./default/go-client/cmd \
  && go run client.go
```

### Protobuf Json (暂时禁用)

Protobuf Json的例子放在`protobufjson`文件夹中。首先需要根据proto文件生成结构体定义。（注：`user.api.go`已经生成，这是为了CI的集成测试使用的，但是我们仍然强烈建议你自己生成一次。）

```shell
cd ./protobufjson \
  && protoc --go_out=. user.proto
```

拷贝`user.api.go`文件到提供者文件夹。

```shell
mv ./protobufjson/user.api.go ./protobufjson/go-server/pkg
```

通过下面的代码启动提供者。

```shell
cd ./protobufjson/go-server/cmd \
  && go run server.go
```

运行消费者来发起泛化调用。

```shell
cd ./protobufjson/go-client/cmd \
  && go run client.go
```
