# Dubbo-go LLM 示例

## 1. **介绍**

本案例展示了如何在 **Dubbo-go** 中集成 **大语言模型（LLM）**，以便在服务端调用 Ollama 模型进行推理，并将结果通过 Dubbo 的 RPC 接口返回给客户端。支持多模型部署和每个模型运行多个实例。

## 2. **准备工作**

### **安装 Ollama**

Ollama 是一个本地运行的大语言模型平台，支持快速推理。

**快速安装**：

```shell
$ curl -fsSL https://ollama.com/install.sh | sh
```

**手动安装**：

```shell
$ mkdir -p ~/ollama
$ cd ~/ollama
$ curl -L https://ollama.com/download/ollama-linux-amd64.tgz -o ollama-linux-amd64.tgz
$ tar -xzf ollama-linux-amd64.tgz
$ echo 'export PATH=$HOME/ollama/bin:$PATH' >> ~/.bashrc
$ source ~/.bashrc
$ ollama serve
```

### 下载模型

```shell
$ ollama pull llava:7b
$ ollama pull qwen2.5:7b  # 可选：下载其他模型
```

您可以自行 pull 需要的模型，并在 `.env` 文件中配置要使用的模型列表。

### **安装 Nacos**

按照此说明[安装并运行 Nacos](https://dubbo-next.staged.apache.org/zh-cn/overview/reference/integrations/nacos/).

## **3. 运行示例**

以下所有的命令都需要在 `llm` 目录下运行。

```shell
$ cd llm
```

生成你的本地配置 `.env` 文件。完成后，请根据实际情况编辑 `.env` 文件并设置相关参数。

```shell
# 复制环境模板文件（Windows用户可使用copy命令）
$ cp .env.example .env
```

### **配置说明**

`.env` 文件支持配置多个模型，示例：

```text
# 支持配置多个模型，使用逗号分隔，支持带空格
OLLAMA_MODELS = llava:7b, qwen2.5:7b
OLLAMA_URL = http://localhost:11434
NACOS_URL = nacos://localhost:8848
TIME_OUT_SECOND = 300
MAX_CONTEXT_COUNT = 3
```

### **服务端运行**

服务端支持多实例部署，每个模型可以运行多个实例以提高服务能力。我们提供了便捷的启动脚本：

**Linux/macOS**:
```shell
# 默认配置：每个模型运行2个实例，端口从20020开始
$ ./start_servers.sh

# 自定义配置：指定实例数量和起始端口
$ ./start_servers.sh --instances 3 --start-port 20030
```

**Windows**:
```shell
# 默认配置：每个模型运行2个实例，端口从20020开始
$ start_servers.bat

# 自定义配置：指定实例数量和起始端口
$ start_servers.bat --instances 3 --start-port 20030
```

### **客户端运行**

客户端通过 Dubbo RPC 调用服务端的接口，获取 Ollama 模型的推理结果。

命令行客户端：
```shell
$ go run go-client/cmd/client.go
```
支持多轮对话、命令交互、上下文管理功能。

Web 客户端：
```shell
$ go run go-client/frontend/main.go
```
访问 `localhost:8080` 使用 Web 界面，支持：
- 多轮对话
- 图片上传（支持 png、jpeg、gif）
- 多模型选择

### **注意事项**

1. 默认超时时间为5分钟（可在 `.env` 中通过 `TIME_OUT_SECOND` 调整）
2. 每个模型默认运行2个实例，可通过启动脚本参数调整
3. 服务端会自动注册到 Nacos，无需手动指定端口
4. 确保所有配置的模型都已通过 Ollama 下载完成