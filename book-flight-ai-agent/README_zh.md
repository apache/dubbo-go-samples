## 订机票示例

### 1. 介绍

本案例展示了 Agent 如何在大语言模型的加持下，完成机票预订的过程。

### 2. 准备工作

#### 修改配置文件

修改配置文件: 复制 `book-flight-ai-agent/.env.example` 为 `book-flight-ai-agent/.env`.

```ini
# LLM 设置
LLM_MODEL = qwq                     # Ollama 模型名称
LLM_URL = "http://127.0.0.1:11434"  # Ollama 的 URL，填写 Ollama 的服务地址
LLM_API_KEY = "sk-..."              # API key

# Client 设置
CLIENT_HOST = "tri://127.0.0.1"     # 客户端主机
CLIENT_PORT = 20000                 # 客户端端口

# Web 设置
WEB_PORT = 8080
TIMEOUT_SECONDS = 300               # 超时时间
```

**注意**：目前仅 Ollama 方式部署的模型

### 3. 运行示例

首先，进入 `book-flight-ai-agent` 目录.

```shell
$ cd book-flight-ai-agent
```

#### 服务端运行

在服务端中集成 Ollama 模型，并使用 Dubbo-go 提供的 RPC 服务进行调用。

在服务端目录下运行：

```shell
$ go run go-server/cmd/server.go
```

#### 客户端运行

前端页面基于 Gin 框架的客户端进行交互，运行以下命令然后访问 ```localhost:8080``` 即可使用:

```shell
$ go run go-client/frntend/main.go
```

### **注意事项**

默认 `Record` 超时时间为两分钟，请确保您的电脑性能能在两分钟内生成相应的响应，否则会超时报错，您也可以在 ```.env``` 文件中自行设置超时时间。
