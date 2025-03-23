## 订机票示例

### 1. 介绍

本案例展示了Agent如何在大语言模型的加持下，完成机票预订的过程。

### 2. 准备工作

#### 修改配置文件

修改配置文件复制`llm/book-flight/go-server/conf/config.example.yml`为`llm/book-flight/go-server/conf/config.example.yml`。

```yaml
LLM: {
  url: "http://localhost:11434",  # Ollama 的 URL，填写 Ollama 的服务地址
  model: "gemma3:27b"             # Ollama 模型名称
}
```

**注意**：目前仅Ollama方式部署的模型

### 3. 运行示例

首先，进入`llm/book-flight`目录.

```shell
$ cd llm/book-flight
```

#### 运行示例

执行下列命令，运行示例

```shell
go run main.go
```

