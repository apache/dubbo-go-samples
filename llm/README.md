# **Dubbo-go LLM Sample**

## 1. **Introduction**

This sample demonstrates how to integrate **large language models (LLM)** in **Dubbo-go**, allowing the server to invoke various LLM providers for inference and return the results to the client via Dubbo RPC. It supports multiple LLM providers including Ollama, OpenAI, Anthropic Claude, and Azure OpenAI, with multiple model deployment and multiple instances per model.

## 2. **Preparation**

### **Choose Your LLM Provider**

This sample supports multiple LLM providers. Choose one based on your needs:

#### **Option 1: Ollama (Local, Free)**
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

**Download Models**:
```shell
$ ollama pull llava:7b
$ ollama pull qwen2.5:7b  # Optional: download additional models
```

#### **Option 2: OpenAI API (Cloud, Paid)**
Requires an OpenAI API key. Get one from [OpenAI Platform](https://platform.openai.com/).

#### **Option 3: Anthropic Claude API (Cloud, Paid)**
Requires an Anthropic API key. Get one from [Anthropic Console](https://console.anthropic.com/).

#### **Option 4: Azure OpenAI (Cloud, Paid)**
Requires an Azure OpenAI resource. Set up at [Azure Portal](https://portal.azure.com/).

### **Install Nacos**

Nacos is required as the service registry for Dubbo-go services. You can install and start Nacos using Docker (recommended) or manually.

**Docker Installation (Recommended)**:

```shell
# Start Nacos in standalone mode
$ docker run -d --name nacos -p 8848:8848 -p 9848:9848 -e MODE=standalone nacos/nacos-server:v2.2.3

# Verify Nacos is running
$ curl http://localhost:8848/nacos/v1/ns/instance/list?serviceName=nacos.naming.service
```

**Manual Installation**:

Follow this instruction to [install and start Nacos server](https://dubbo-next.staged.apache.org/zh-cn/overview/reference/integrations/nacos/).

**Verify Nacos Installation**:

```shell
# Check if Nacos is accessible
$ curl http://localhost:8848/nacos/v1/ns/instance/list?serviceName=nacos.naming.service

# Expected response: {"name":"DEFAULT_GROUP@@nacos.naming.service",...}
```

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

The `.env` file supports multiple LLM providers. Here are examples for each provider:

#### **Ollama Configuration**:
```text
LLM_PROVIDER = ollama
LLM_MODELS = llava:7b, qwen2.5:7b
LLM_BASE_URL = http://localhost:11434
MODEL_NAME = llava:7b
NACOS_URL = nacos://localhost:8848
TIME_OUT_SECOND = 300
MAX_CONTEXT_COUNT = 3
```

#### **OpenAI Configuration**:
```text
LLM_PROVIDER = openai
LLM_MODELS = gpt-4, gpt-3.5-turbo
LLM_BASE_URL = https://api.openai.com/v1
LLM_API_KEY = your-openai-api-key
MODEL_NAME = gpt-4
NACOS_URL = nacos://localhost:8848
TIME_OUT_SECOND = 300
MAX_CONTEXT_COUNT = 3
```

#### **Anthropic Configuration**:
```text
LLM_PROVIDER = anthropic
LLM_MODELS = claude-3-sonnet-20240229, claude-3-haiku-20240307
LLM_BASE_URL = https://api.anthropic.com/v1
LLM_API_KEY = your-anthropic-api-key
MODEL_NAME = claude-3-sonnet-20240229
NACOS_URL = nacos://localhost:8848
TIME_OUT_SECOND = 300
MAX_CONTEXT_COUNT = 3
```

#### **Azure OpenAI Configuration**:
```text
LLM_PROVIDER = azure-openai
LLM_MODELS = gpt-4, gpt-35-turbo
LLM_BASE_URL = https://your-resource.openai.azure.com/openai/deployments/your-deployment
LLM_API_KEY = your-azure-openai-api-key
MODEL_NAME = gpt-4
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