## Dubbo-Go Group Usage

### 1. Introduction

When you have multi-impls of a interface,you can distinguish them with the group.

### 2. How to configure the group

provider side

```yaml
# server.yml

# service config
services:
  "UserProvider":
    # ...
    group: "group-a"
    # ...
```

consumer side

```yaml
# client.yml

# reference config
references:
  "UserProvider":
    # ...
    group: "GroupA"
    # ...
```