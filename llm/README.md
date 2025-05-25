# **Dubbo-go LLM Sample**

## 1. **Introduction**

This sample demonstrates how to integrate **large language models (LLM)** in **Dubbo-go**, allowing the server to invoke the Ollama model for inference and return the results to the client via Dubbo RPC. It supports multiple model deployment with multiple instances per model.

## 2. **Preparation**

### **Install Ollama**

Ollama is a local language model platform that supports fast inference.

**Quick Installation**:

```shell
$ curl -fsSL https://ollama.com/install.sh | sh
```

**Manual Installation**:

```shell
$ mkdir -p ~/ollama
$ cd ~/ollama
$ curl -L https://ollama.com/download/ollama-linux-amd64.tgz -o ollama-linux-amd64.tgz
$ tar -xzf ollama-linux-amd64.tgz
$ echo 'export PATH=$HOME/ollama/bin:$PATH' >> ~/.bashrc
$ source ~/.bashrc
$ ollama serve
```

### **Download Models**

```shell
$ ollama pull llava:7b
$ ollama pull qwen2.5:7b  # Optional: download additional models
```

You can pull your preferred models and configure them in the `.env` file.

### **Install Nacos**

Follow this instruction to [install and start Nacos server](https://dubbo-next.staged.apache.org/zh-cn/overview/reference/integrations/nacos/).

## 3. **Run the Example**

You need to run all the commands in the `llm` directory.

```shell
$ cd llm
```

Create your local environment configuration by copying the template file. 
After creating the `.env` file, edit it to set up your specific configurations.

```shell
# Copy environment template (Use `copy` for Windows)
$ cp .env.example .env
```

### **Configuration**

The `.env` file supports multiple model configurations, example:

```text
# Configure multiple models, comma-separated, spaces allowed
OLLAMA_MODELS = llava:7b, qwen2.5:7b
OLLAMA_URL = http://localhost:11434
NACOS_URL = nacos://localhost:8848
TIME_OUT_SECOND = 300
MAX_CONTEXT_COUNT = 3
```

### **Run the Server**

The server supports multi-instance deployment, with multiple instances per model to enhance service capacity. We provide convenient startup scripts:

**Linux/macOS**:
```shell
# Default: 2 instances per model, starting from port 20020
$ ./start_servers.sh

# Custom configuration: specify instance count and start port
$ ./start_servers.sh --instances 3 --start-port 20030
```

**Windows**:
```shell
# Default: 2 instances per model, starting from port 20020
$ start_servers.bat

# Custom configuration: specify instance count and start port
$ start_servers.bat --instances 3 --start-port 20030
```

### **Run the Client**

The client invokes the server's RPC interface to retrieve inference results from the Ollama models.

CLI Client:
```shell
$ go run go-client/cmd/client.go
```
Supports multi-turn conversations, command interaction, and context management.

Web Client:
```shell
$ go run go-client/frontend/main.go
```
Access at `localhost:8080` with features:
- Multi-turn conversations
- Image upload support (png, jpeg, gif)
- Multiple model selection

### **Important Notes**

1. Default timeout is 5 minutes (adjustable via `TIME_OUT_SECOND` in `.env`)
2. Each model runs 2 instances by default, adjustable via startup script parameters
3. Servers automatically register with Nacos, no manual port specification needed
4. Ensure all configured models are downloaded through Ollama before starting