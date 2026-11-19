# Triple compatibility and Java interoperability

From the repository root, run:

```bash
bash ./integrate_test.sh rpc/triple/compatibility
```

The same scenario is included in `start_integrate_test.sh`. It runs:

1. A dubbo-go `v3.3.2` client against the server built with this repository's selected dubbo-go dependency.
2. The current client against a separately built dubbo-go `v3.3.2` server.
3. The existing `java_interop/non-protobuf-triple` sample in all four Go/Java client-server combinations.

All calls use Hessian2 over Triple. Both cross-version calls assert the returned greeting, and the Java matrix checks the expected response from each direction. The existing Java sample pins Dubbo `3.3.0-beta.3`.

Go, Java, Maven, Make, Python 3, and network access for dependencies are required. TCP port `20000` must be free. Binaries and matrix logs are created in a temporary directory; the script cleans up its servers and temporary files on exit.

The current client and server use the root module's dubbo-go dependency, including any replacement configured by the caller or CI. The `legacy-client` nested module deliberately pins `v3.3.2` without inheriting that replacement: Go cannot load two versions of the same module path in one build. Root `go test ./...` does not run this integration scenario.
