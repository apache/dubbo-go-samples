# **Dubbo-go LLM Sample**

## 1. **Introduction**

This sample demonstrates how to integrate **large language models (LLM)** in **Dubbo-go**, allowing the server to invoke the Ollama model for inference and return the results to the client via Dubbo RPC.

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

### **Download Model**

```shell
$ ollama run deepseek-r1:1.5b
```

## 3. **Run the Example**

You need to run all the commands in ```llm``` directory.

```shell
$ cd llm
```
### **Run the Server**

The server integrates the Ollama model and uses Dubbo-go's RPC service for invocation.

Run the server by executing:

```shell
$ go run go-server/cmd/server.go
```

### **Run the Client**

The client invokes the server's RPC interface to retrieve the inference results from the Ollama model.

Run the cli client by executing:

```shell
$ go run go-client/cmd/client.go
```

Cli client supports multi-turn conversations, command interact, context management.

We also support a frontend using Gin framework for users to interact. If you want run the frontend client you can executing the following command and open it in ```localhost:8080``` by default:

```shell
$ go run go-client/frontend/main.go
```

Frontend client supports multi-turn conversations, binary file (image) support for LLM interactions.
Currently the supported uploaded image types are limited to png, jpeg and gif, with plans to support more binary file types in the future.
