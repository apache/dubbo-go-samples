#!/bin/bash

# Check if .env file exists
if [ ! -f .env ]; then
    echo "Error: .env file not found"
    exit 1
fi

# Default values
START_PORT=20020
INSTANCES_PER_MODEL=2

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --instances)
            INSTANCES_PER_MODEL=$2
            shift 2
            ;;
        --start-port)
            START_PORT=$2
            shift 2
            ;;
        *)
            echo "Unknown parameter: $1"
            exit 1
            ;;
    esac
done

# Function to read value from .env file and clean it
get_env_value() {
    local key=$1
    # Read the line containing the key, extract everything after the first =
    local value=$(grep "^[[:space:]]*$key[[:space:]]*=" .env | sed 's/^[^=]*=[[:space:]]*//')
    # Remove leading/trailing whitespace and quotes
    value=$(echo "$value" | sed -e 's/^[[:space:]]*//' -e 's/[[:space:]]*$//' -e 's/^["\x27]//' -e 's/["\x27]$//')
    echo "$value"
}

# Get models from .env file
OLLAMA_MODELS=$(get_env_value "OLLAMA_MODELS")

# Check if OLLAMA_MODELS is empty
if [ -z "$OLLAMA_MODELS" ]; then
    echo "Error: OLLAMA_MODELS not found in .env file"
    echo "Please make sure .env file contains a line like: OLLAMA_MODELS = llava:7b, qwen2.5:7b"
    exit 1
fi

echo "Found models: $OLLAMA_MODELS"

# Convert comma-separated string to array, handling spaces
IFS=',' read -ra MODELS <<< "$OLLAMA_MODELS"

current_port=$START_PORT

# Function to start a model instance
start_model_instance() {
    local model=$1
    local port=$2
    local instance_num=$3
    
    # Clean the model name (remove leading/trailing spaces)
    model=$(echo "$model" | sed -e 's/^[[:space:]]*//' -e 's/[[:space:]]*$//')
    
    echo "Starting server for model: $model (Instance $instance_num, Port: $port)"
    export MODEL_NAME=$model
    export SERVER_PORT=$port
    go run go-server/cmd/server.go &
    sleep 2
}

# Start instances for each model
for model in "${MODELS[@]}"; do
    echo "Processing model: $model"
    for ((i=1; i<=INSTANCES_PER_MODEL; i++)); do
        start_model_instance "$model" $current_port $i
        ((current_port++))
    done
done

echo "All servers started. Total instances: $((${#MODELS[@]} * INSTANCES_PER_MODEL))"
echo "Use Ctrl+C to stop all servers."

# Wait for all background processes
wait 