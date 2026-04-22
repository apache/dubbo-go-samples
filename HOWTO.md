# How To Run

English | [中文](HOWTO_CN.md)

There are two ways to run dubbo-go samples:

1. Script-driven integration testing (CI flow)
2. Manual run without scripts

## 1. Integration Testing (CI flow)

### 1.1 Integration flow overview

Current CI uses script-driven integration tests.

For one sample, `./integrate_test.sh <sample-path>` executes:
1. Start `go-server`
2. Start auxiliary Go servers (`*-server/cmd/*.go`, e.g. `grpc-server`)
3. Run `go-client`
4. Run `java-client` (if present)
5. Stop `go-server`
6. Start `java-server` (if present)
7. Run `java-client` (if present)
8. Run `go-client` again (verify Go client can call Java server)

Notes:
- If `mvn` is unavailable, Java phases are skipped automatically.
- Java server readiness is checked by TCP port before Java/Go client phases continue.
- Default Java server endpoint is `127.0.0.1:20000` (can be overridden by env vars in `integrate_test.sh`).

### 1.2 How to use `start_integrate_test.sh` and `integrate_test.sh`

Prerequisite:
- Docker / Docker Compose
- Go toolchain
- Maven (optional but required for Java phases)

Run all samples in CI list:
```bash
./start_integrate_test.sh
```

`start_integrate_test.sh` will:
- Start dependencies via root `docker-compose.yml`
- Health-check dependencies
- Run each sample by calling `./integrate_test.sh <sample>`
- Stop dependencies at exit (success or failure)

Run a single sample:
```bash
./integrate_test.sh helloworld
./integrate_test.sh direct
```

Useful env vars:
- `GO_CLIENT_TIMEOUT_SECONDS` (default: `90`)
- `JAVA_CLIENT_TIMEOUT_SECONDS` (default: `180`)
- `JAVA_SERVER_READY_TIMEOUT_SECONDS` (default: `60`)
- `JAVA_SERVER_HOST` (default: `127.0.0.1`)
- `JAVA_SERVER_PORT` (default: `20000`)

### 1.3 How to add a new integration test sample

Create a sample directory with at least:
- `go-server/cmd/*.go`
- `go-client/cmd/*.go`

- Optional Java interop:
  - `java-server/run.sh`
  - `java-client/run.sh`

Requirements for Java scripts:
- `java-server/run.sh` should keep server process alive (do not rely on stdin in background mode).
- `java-client/run.sh` should exit with non-zero on failure.

Validation steps:
1. Run the sample only:
   ```bash
   ./integrate_test.sh <your-sample-path>
   ```
2. Add it into the `array` list in `start_integrate_test.sh`.
3. Run full flow:
   ```bash
   ./start_integrate_test.sh
   ```
4. Ensure failure paths are visible (non-zero exit code, useful logs).

## 2. Manual run (without scripts)

This section shows how to run one sample manually, without `start_integrate_test.sh` or `integrate_test.sh`.

Example sample: `helloworld`

### 2.1 Start dependencies

```bash
cd <PATH OF dubbo-go-samples>
docker compose -f docker-compose.yml up -d
```

If your environment uses legacy compose:
```bash
docker-compose -f docker-compose.yml up -d
```

### 2.2 Run Go server

Open a new terminal:
```bash
cd <PATH OF dubbo-go-samples>/helloworld
export DUBBO_GO_CONFIG_PATH=./go-server/conf/dubbogo.yml
go run ./go-server/cmd/*.go
```

### 2.3 Run Go client

Open another terminal:
```bash
cd <PATH OF dubbo-go-samples>/helloworld
export DUBBO_GO_CONFIG_PATH=./go-client/conf/dubbogo.yml
go run ./go-client/cmd/*.go
```

### 2.4 Optional Java interop verification

1. Stop Go server (Ctrl+C in Go server terminal).
2. Start Java server:
   ```bash
   cd <PATH OF dubbo-go-samples>/helloworld/java-server
   bash ./run.sh
   ```
3. Run Java client:
   ```bash
   cd <PATH OF dubbo-go-samples>/helloworld/java-client
   bash ./run.sh
   ```
4. Run Go client again (from step 2.3 terminal) to verify Go -> Java call.

### 2.5 Cleanup

Stop foreground processes with `Ctrl+C`, then:
```bash
cd <PATH OF dubbo-go-samples>
docker compose -f docker-compose.yml down
```
