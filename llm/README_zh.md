# Dubbo-go LLM 示例

## 1. **介绍**

本案例展示了如何在 **Dubbo-go** 中集成 **大语言模型（LLM）**，以便在服务端调用 Ollama 模型进行推理，并将结果通过 Dubbo 的 RPC 接口返回给客户端。

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
$ ollama run deepseek-r1:1.5b
```

## **3. 运行示例**

以下所有的命令都需要在```llm``` 目录下运行.

```shell
$ cd llm
```
### **服务端运行**

在服务端中集成 Ollama 模型，并使用 Dubbo-go 提供的 RPC 服务进行调用。

在服务端目录下运行：

```shell
$ go run go-server/cmd/server.go
```

### **客户端运行**

客户端通过 Dubbo RPC 调用服务端的接口，获取 Ollama 模型的推理结果。

在客户端目录下运行：

```shell
$ go run go-client/cmd/client.go
```

命令行客户端支持多轮对话、命令交互、上下文管理功能。

我们也提供了包含前端页面的基于Gin框架的客户端进行交互，运行以下命令然后访问 ```localhost:8080``` 即可使用:

```shell
$ go run go-client/frontend/main.go
```

包含前端页面的客户端支持多轮对话，支持进行二进制文件（图片）传输并与大模型进行交互。
目前所支持上传的图片类型被限制为 png，jpeg 和 gif 类型，计划在将来支持更多的二进制文件类型。