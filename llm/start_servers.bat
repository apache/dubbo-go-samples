@echo off
setlocal enabledelayedexpansion

REM Default values
set START_PORT=20020
set INSTANCES_PER_MODEL=2

REM Parse command line arguments
:parse_args
if "%~1"=="" goto :end_parse
if "%~1"=="--instances" (
    set INSTANCES_PER_MODEL=%~2
    shift
    shift
    goto :parse_args
)
if "%~1"=="--start-port" (
    set START_PORT=%~2
    shift
    shift
    goto :parse_args
)
:end_parse

REM Check if .env file exists
if not exist .env (
    echo Error: .env file not found
    exit /b 1
)

REM Read and clean OLLAMA_MODELS from .env file
for /f "usebackq tokens=*" %%a in (`powershell -Command "Get-Content .env | Select-String '^[[:space:]]*OLLAMA_MODELS[[:space:]]*=' | ForEach-Object { $_.Line -replace '^[^=]*=[[:space:]]*', '' -replace '^[\"'']+|[\"'']+$', '' }"`) do (
    set MODELS_LINE=%%a
)

if "!MODELS_LINE!"=="" (
    echo Error: OLLAMA_MODELS not found in .env file
    echo Please make sure .env file contains a line like: OLLAMA_MODELS = llava:7b, qwen2.5:7b
    exit /b 1
)

echo Found models: !MODELS_LINE!

REM Split models and start servers
set current_port=%START_PORT%

for %%m in (!MODELS_LINE:,=^,!) do (
    set "model=%%m"
    set "model=!model: =!"
    echo Processing model: !model!
    
    for /l %%i in (1,1,%INSTANCES_PER_MODEL%) do (
        echo Starting server for model: !model! (Instance %%i, Port: !current_port!)
        set MODEL_NAME=!model!
        set SERVER_PORT=!current_port!
        start /B cmd /c "go run go-server/cmd/server.go"
        timeout /t 2 /nobreak >nul
        set /a current_port+=1
    )
)

set /a total_instances=0
for %%a in (!MODELS_LINE:,=^,!) do set /a total_instances+=1
set /a total_instances*=%INSTANCES_PER_MODEL%

echo All servers started. Total instances: !total_instances!
echo Use Ctrl+C to stop all servers.
pause 