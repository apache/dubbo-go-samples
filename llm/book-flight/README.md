## Book Flight example

### 1. Introduction

This case shows how Agent completes the airline booking process with the support of a large language model.

### 2. Preparation

#### Modify the configuration file

Modify the configuration file and copy `llm/book-flight/go-server/conf/config.example.yml` to `llm/book-flight/go-server/conf/config.example.yml`.

```yaml
LLM: {
url: "http://localhost:11434", # Ollama URL, fill in Ollama service address
model: "gemma3:27b" # Ollama model name
}
```

**Note**: Currently only models deployed in Ollama mode

### 3. Run the example

First, enter the `llm/book-flight` directory.

```shell
$ cd llm/book-flight
```

#### Run the example

Execute the following command to run the example

```shell
go run main.go
```