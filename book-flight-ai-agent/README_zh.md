## 订机票示例

### 1. 介绍

本案例展示了 Agent 如何在大语言模型的加持下，完成机票预订的过程。

### 2. 准备工作

#### 修改配置文件

修改配置文件: 复制 `book-flight-ai-agent/.env.example` 为 `book-flight-ai-agent/.env`.

```ini
# LLM 设置
LLM_MODEL = qwq                     # Ollama 模型名称（详见下方"模型选择建议"）
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

#### 模型选择建议

默认模型是 **`qwq`**。本示例的设计目标是演示**多模态 Agent**：聊天协议带了可选的图片字段（`proto/chat.proto` 里的 `bytes bin`），前端也允许在消息里附图。`qwq` 是当时 Ollama 上能找到的、**同时满足 ReAct 质量和多模态输入**的最小模型；它在 Ollama 中的实际显存/内存占用对多数消费级机器并不算重。

如果你**不需要多模态能力**，只想用纯文本查询机票，可以换成更小的纯文本模型。但要注意：Agent 的工具选择对模型的指令遵循能力要求较高，模型过小会导致请求被误判为 `TaskUnrelated`，agent 根本不调用查询工具。

| 使用场景 | 推荐模型 | 大致占用 |
|---------|---------|---------|
| **多模态**（默认，推荐） | `qwq` | Ollama 下并不重，普通显卡/CPU 都能跑 |
| 纯文本，多轮推理质量优先 | `qwen2.5:14b` | ~8 GB |
| 纯文本，机器配置较低 | `qwen2.5:7b` | ~4 GB |

启动服务端前先拉模型：

```shell
$ ollama pull qwq          # 默认
# 或者，做纯文本测试时：
$ ollama pull qwen2.5:7b
```

如果本地完全跑不了 LLM，但又想验证查询/购买工具是否正确，可以直接跑单元测试：

```shell
$ go test ./go-server/tools/bookingflight/... -v
```

> **提示**：小于约 7B 的模型（例如 `qwen3:4b`）经常把用户意图识别成 `TaskUnrelated`，导致 agent 根本不调查询工具。这是模型指令遵循能力不足，跟代码无关。看到 agent 不调任何工具时，先试试换更大的模型，再去查代码。

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
$ go run go-client/frontend/main.go
```

### **注意事项**

默认 `Record` 超时时间为两分钟，请确保您的电脑性能能在两分钟内生成相应的响应，否则会超时报错，您也可以在 ```.env``` 文件中自行设置超时时间。
