## Dubbo-Go 服务分组

### 1. 介绍

当一个接口有多种实现时，可以用 group 区分。

### 2. 使用示例

生产者

```yaml
# server.yml

# service config
services:
  "UserProvider":
    # ...
    group: "group-a"
    # ...
```

消费者

```yaml
# client.yml

# reference config
references:
  "UserProvider":
    # ...
    group: "GroupA"
    # ...
```