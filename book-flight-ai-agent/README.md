## Book Flight example

### 1. Introduction

This case shows how Agent completes the airline booking process with the support of a large language model.

### 2. Preparation

#### Modify the configuration file

Modify the configuration file and copy `book-flight-ai-agent/.env.example` to `book-flight-ai-agent/.env`.

```ini
# LLM Settings
LLM_MODEL = qwq # Ollama model name (see "Model Recommendations" below)
LLM_URL = "http://127.0.0.1:11434" # Ollama URL, fill in Ollama service address
LLM_API_KEY = "sk-..." # API key

# Client Settings
CLIENT_HOST = "tri://127.0.0.1"
CLIENT_PORT = 20000

# Web Settings
WEB_PORT = 8080
TIMEOUT_SECOND = 300 # Timeout
```

**Note**: Currently only models deployed in Ollama mode

#### Model Recommendations

The default model is **`qwq`**. This sample is designed to demonstrate a multimodal agent — the chat protocol carries an optional image (`bytes bin` in `proto/chat.proto`), and the front-end allows uploading a picture together with the message. `qwq` was chosen as the default because, at the time the sample was written, it was the smallest Ollama-distributed model that combined ReAct-quality reasoning with multimodal input. In practice, its actual VRAM/RAM usage in Ollama is modest for many consumer machines.

If you do **not** need the multimodal capability and just want to exercise the agent on text-only flight queries, you can swap to a smaller text-only model. Be aware that the agent's tool selection still requires solid instruction-following, so going too small will make the model misclassify the request as `TaskUnrelated` and never invoke the search tool.

| Use case | Suggested model | Approx. footprint |
|----------|------------------|-------------------|
| **Multimodal** (default, recommended) | `qwq` | Modest in Ollama, runs fine on consumer GPUs/CPU |
| Text-only, multi-turn quality | `qwen2.5:14b` | ~8 GB |
| Text-only, lighter machine | `qwen2.5:7b` | ~4 GB |

Pull the model before starting the server:

```shell
$ ollama pull qwq          # default
# or, for text-only testing:
$ ollama pull qwen2.5:7b
```

If you can't run any LLM locally and just want to verify the booking tools, run the unit tests directly:

```shell
$ go test ./go-server/tools/bookingflight/... -v
```

> **Tip**: Models smaller than ~7B (e.g. `qwen3:4b`) tend to misclassify the user's intent and route to `TaskUnrelated`, so the agent never invokes the search tool. The code itself is fine — the issue is purely the model's instruction-following ability. If you see the agent refuse to call any tool, try a larger model first before debugging the code.

### 3. Run the example

First, enter the `book-flight-ai-agent` directory.

```shell
$ cd book-flight-ai-agent
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
$ go run go-client/frontend/main.go
```

### **Notes**

The default `Record` timeout is two minutes. Please ensure that your computer performance can generate the corresponding response within two minutes, otherwise it will time out and report an error. You can also set the timeout in the ```.env``` file.