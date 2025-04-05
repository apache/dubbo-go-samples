## Book Flight example

### 1. Introduction

This case shows how Agent completes the airline booking process with the support of a large language model.

### 2. Preparation

#### Modify the configuration file

Modify the configuration file and copy `llm/book-flight/.env.example` to `llm/book-flight/.env`.

```ini
LLM_MODEL = qwq # Ollama model name
LLM_URL = "http://127.0.0.1:11434" # Ollama URL, fill in Ollama service address
LLM_API_KEY = "sk-..." # API key
TIME_OUT_SECOND = 300 # Timeout
```

**Note**: Currently only models deployed in Ollama mode

### 3. Run the example

First, enter the `llm/book-flight` directory.

```shell
$ cd llm/book-flight
```

#### Server operation

Integrate the Ollama model in the server and call it using the RPC service provided by Dubbo-go.

Run in the server directory:

```shell
$ go run go-server/cmd/server.go
```

#### Client operation

The front-end page interacts with the client based on the Gin framework. Run the following command and then visit ```localhost:8080``` to use it:

```shell
$ go run go-client/frntend/main.go
```

### **Notes**

The default `Record` timeout is two minutes. Please ensure that your computer performance can generate the corresponding response within two minutes, otherwise it will time out and report an error. You can also set the timeout in the ```.env``` file.